package shared

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"net/http/httptrace"
	"net/url"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/imroc/req/v3"
)

const defaultMaxRetries = 5

const retryWindow = 5 * time.Minute

type Clock interface {
	Now() time.Time
	Sleep(ctx context.Context, d time.Duration) error
}

type systemClock struct{}

func (systemClock) Now() time.Time {
	return time.Now()
}

func (systemClock) Sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

var SystemClock Clock = systemClock{}

type Retries struct {
	clock       Clock
	maxRetries  atomic.Int64
	backoff     func(attempt int) time.Duration
	keysHonored atomic.Bool
}

func NewRetries(clock Clock) *Retries {
	retries := &Retries{clock: clock, backoff: backoff(400*time.Millisecond, 10*time.Second)}
	retries.maxRetries.Store(defaultMaxRetries)
	return retries
}

func (s *Retries) SetMaxRetries(client *req.Client, count int) {
	s.maxRetries.Store(int64(count))
	client.SetCommonRetryCount(count)
}

func (s *Retries) remember(status idempotencyKeyStatus) {
	if status.known() {
		s.keysHonored.Store(status.keyHonored())
	}
}

func (s *Retries) honorsKeys() bool {
	return s.keysHonored.Load()
}

func (s *Retries) interval(resp *req.Response, attempt int) time.Duration {
	ctx := context.Background()
	if resp != nil && resp.Request != nil {
		ctx = resp.Request.Context()
	}
	_ = s.clock.Sleep(ctx, s.backoff(attempt))
	return 0
}

type writeContextKey struct{}

type write struct {
	retries                 *Retries
	start                   time.Time
	mayHaveReachedTheServer bool
	mayHaveRun              bool
}

func writeOf(resp *req.Response) *write {
	if resp == nil || resp.Request == nil {
		return nil
	}
	w, _ := resp.Request.Context().Value(writeContextKey{}).(*write)
	return w
}

func (s *Retries) sendWrites(rt req.RoundTripper) req.RoundTripFunc {
	return func(r *req.Request) (*req.Response, error) {
		if isSafeMethod(r.Method) {
			return rt.RoundTrip(r)
		}
		if err := setIdempotencyKey(r); err != nil {
			return &req.Response{Request: r, Err: err}, err
		}
		w := &write{retries: s, start: s.clock.Now()}
		ctx := context.WithValue(r.Context(), writeContextKey{}, w)
		defer r.SetContext(ctx)
		for retryCount := 0; ; retryCount++ {
			var gotConn atomic.Bool
			r.SetContext(httptrace.WithClientTrace(ctx, &httptrace.ClientTrace{
				GotConn: func(httptrace.GotConnInfo) { gotConn.Store(true) },
			}))
			resp, err := rt.RoundTrip(r)
			if !w.attempted(resp, err, gotConn.Load()) || !s.waitBeforeRetry(ctx, retryCount, w) {
				return resp, err
			}
		}
	}
}

func (w *write) attempted(resp *req.Response, err error, gotConn bool) bool {
	if resp != nil && resp.Response != nil {
		w.mayHaveReachedTheServer = true
		w.retries.remember(idempotencyKeyStatusOf(resp.Header))
		if err != nil {
			if !isResponseLost(resp, err) {
				return false
			}
			w.mayHaveRun = true
			return w.honorsKeys()
		}
		if isKeyInUse(resp) {
			w.mayHaveRun = true
			return true
		}
		return resp.StatusCode == 429
	}
	switch classifyFailure(err, gotConn) {
	case neverSent:
		return true
	case responseLost:
		w.mayHaveReachedTheServer, w.mayHaveRun = true, true
		return w.honorsKeys()
	default:
		w.mayHaveReachedTheServer, w.mayHaveRun = true, true
		return false
	}
}

func (w *write) honorsKeys() bool {
	return w.retries.honorsKeys()
}

func (w *write) startsWithinRetryWindow(delay time.Duration) bool {
	return w.retries.clock.Now().Add(delay).Sub(w.start) <= retryWindow
}

func (w *write) failureWithoutResponse() (outcomeUnknown bool, safeToRetryWithSameKey bool) {
	return w.mayHaveRun, !w.mayHaveRun || (w.honorsKeys() && w.startsWithinRetryWindow(0))
}

func (s *Retries) waitBeforeRetry(ctx context.Context, retryCount int, w *write) bool {
	if int64(retryCount) >= s.maxRetries.Load() {
		return false
	}
	delay := s.backoff(retryCount + 1)
	if w.mayHaveReachedTheServer && !w.startsWithinRetryWindow(delay) {
		return false
	}
	if ctx.Err() != nil {
		return false
	}
	if s.clock.Sleep(ctx, delay) != nil {
		return false
	}
	return !w.mayHaveReachedTheServer || w.startsWithinRetryWindow(0)
}

type failureKind int

const (
	neverSent failureKind = iota
	responseLost
	notRetryable
)

func classifyFailure(err error, gotConn bool) failureKind {
	if gotConn {
		return responseLost
	}
	if isNetworkError(err) || isTLSError(err) || errors.Is(err, context.Canceled) {
		return neverSent
	}
	return notRetryable
}

func isResponseLost(resp *req.Response, err error) bool {
	return resp.Response == nil || isNetworkError(err)
}

func isNetworkError(err error) bool {
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		err = urlErr.Err
	}
	var opErr *net.OpError
	var dnsErr *net.DNSError
	var timeout interface{ Timeout() bool }
	return errors.As(err, &opErr) || errors.As(err, &dnsErr) || (errors.As(err, &timeout) && timeout.Timeout()) ||
		errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, syscall.EPIPE)
}

func isTLSError(err error) bool {
	var recordHeaderErr tls.RecordHeaderError
	var verificationErr *tls.CertificateVerificationError
	var alertErr tls.AlertError
	var unknownAuthorityErr x509.UnknownAuthorityError
	var hostnameErr x509.HostnameError
	var invalidErr x509.CertificateInvalidError
	return errors.As(err, &recordHeaderErr) || errors.As(err, &verificationErr) || errors.As(err, &alertErr) ||
		errors.As(err, &unknownAuthorityErr) || errors.As(err, &hostnameErr) || errors.As(err, &invalidErr)
}
