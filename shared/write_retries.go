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
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/imroc/req/v3"
)

const defaultMaxRetries = 5

// retryWindow bounds the retries of a write once an attempt may have reached the server. The server keeps an
// Idempotency-Key for 10 minutes from the first attempt; a retry after that would run as a new request.
const retryWindow = 5 * time.Minute

// Clock is the time source of the retries. It only exists so that tests don't have to wait.
type Clock interface {
	// Now returns the current time. time.Now carries a monotonic reading, so differences don't jump with the wall clock.
	Now() time.Time
	// Sleep waits for d, or until ctx ends, in which case it returns the context's error.
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

// SystemClock is the real clock.
var SystemClock Clock = systemClock{}

// Retries is the retry state of one API client: the maximum number of retries, the clock, and per operation whether the
// server honors idempotency keys for it. Not shared between clients, because that depends on the company of the
// secret key.
type Retries struct {
	clock      Clock
	maxRetries atomic.Int64
	backoff    func(attempt int) time.Duration
	// routes maps "METHOD path template" to whether the last response with a known Idempotency-Key-Status said the
	// route honors keys. No entry means unknown.
	routes sync.Map
}

func NewRetries(clock Clock) *Retries {
	retries := &Retries{clock: clock, backoff: backoff(400*time.Millisecond, 10*time.Second)}
	retries.maxRetries.Store(defaultMaxRetries)
	return retries
}

// SetMaxRetries sets the number of retries after the first attempt, for every reason together. 0 disables retries.
func (s *Retries) SetMaxRetries(client *req.Client, count int) {
	s.maxRetries.Store(int64(count))
	client.SetCommonRetryCount(count)
}

func (s *Retries) remember(route string, status idempotencyKeyStatus) {
	if status.known() {
		s.routes.Store(route, status.keyHonored())
	}
}

func (s *Retries) honorsKeys(route string) bool {
	honored, ok := s.routes.Load(route)
	return ok && honored.(bool)
}

// interval is the wait between two attempts of a GET. It waits itself, and then tells req not to wait any longer: req
// waits with time.Sleep, which a cancelled or expired context cannot interrupt.
func (s *Retries) interval(resp *req.Response, attempt int) time.Duration {
	ctx := context.Background()
	if resp != nil && resp.Request != nil {
		ctx = resp.Request.Context()
	}
	_ = s.clock.Sleep(ctx, s.backoff(attempt))
	return 0
}

type writeContextKey struct{}

// write holds the attempts of one POST, PUT, PATCH or DELETE call, all sent with the same idempotency key.
type write struct {
	retries *Retries
	route   string
	start   time.Time
	// mayHaveReachedTheServer is set once a response arrived or a response was lost. From then on, retries must start
	// within the retry window.
	mayHaveReachedTheServer bool
	// mayHaveRun is set by a lost response, a 409 IDEMPOTENCY_KEY_IN_USE, or a failure that isn't a network failure. It
	// decides whether a later failure to send leaves the outcome unknown.
	mayHaveRun bool
}

func writeOf(resp *req.Response) *write {
	if resp == nil || resp.Request == nil {
		return nil
	}
	w, _ := resp.Request.Context().Value(writeContextKey{}).(*write)
	return w
}

// sendWrites sends every request that isn't a GET, HEAD, OPTIONS or TRACE with an idempotency key, and retries it with
// the same key: after a 429 or a 409 IDEMPOTENCY_KEY_IN_USE, when it was never sent, and when its response was lost on
// a route that the server said honors keys. All retries count against maxRetries. req's own retries are switched off
// for these requests by the retry condition.
func (s *Retries) sendWrites(rt req.RoundTripper) req.RoundTripFunc {
	return func(r *req.Request) (*req.Response, error) {
		if isSafeMethod(r.Method) {
			return rt.RoundTrip(r)
		}
		if err := setIdempotencyKey(r); err != nil {
			return &req.Response{Request: r, Err: err}, err
		}
		w := &write{retries: s, route: r.Method + " " + routeTemplate(r.RawURL), start: s.clock.Now()}
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

// routeTemplate is the path before req fills in the path parameters, e.g. /events/{event}/actions/book, without a
// query string.
func routeTemplate(rawURL string) string {
	if i := strings.IndexByte(rawURL, '?'); i >= 0 {
		return rawURL[:i]
	}
	return rawURL
}

// attempted records the outcome of an attempt and reports whether it may be retried, retries left and time permitting.
func (w *write) attempted(resp *req.Response, err error, gotConn bool) bool {
	if resp != nil && resp.Response != nil {
		w.mayHaveReachedTheServer = true
		w.retries.remember(w.route, idempotencyKeyStatusOf(resp.Header))
		if err != nil {
			if !isResponseLost(resp, err) {
				return false
			}
			w.mayHaveRun = true
			return w.routeHonorsKeys()
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
		return w.routeHonorsKeys()
	default:
		w.mayHaveReachedTheServer, w.mayHaveRun = true, true
		return false
	}
}

func (w *write) routeHonorsKeys() bool {
	return w.retries.honorsKeys(w.route)
}

func (w *write) startsWithinRetryWindow(delay time.Duration) bool {
	return w.retries.clock.Now().Add(delay).Sub(w.start) <= retryWindow
}

// failureWithoutResponse returns the outcome of a write whose last attempt got no (complete) response.
func (w *write) failureWithoutResponse() (outcomeUnknown bool, safeToRetryWithSameKey bool) {
	return w.mayHaveRun, !w.mayHaveRun || (w.routeHonorsKeys() && w.startsWithinRetryWindow(0))
}

// waitBeforeRetry reports whether another attempt may start, after waiting for it. maxRetries is the hard limit and is
// checked first; the retry window can only refuse a retry that maxRetries allows.
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
	return s.clock.Sleep(ctx, delay) == nil
}

type failureKind int

const (
	// neverSent: no connection was obtained, so no byte of the request reached the server (connection refused,
	// connect timeout, DNS failure, TLS handshake failure).
	neverSent failureKind = iota
	// responseLost: the request may have reached the server and run.
	responseLost
	// notRetryable: not a network failure; a retry would fail the same way.
	notRetryable
)

// classifyFailure classifies a failed attempt that got no response. Only an attempt that never got a connection is
// never sent: req's transport reports a connection (GotConn) before it writes the request, for HTTP/1 and HTTP/2.
// When in doubt, the response was lost.
func classifyFailure(err error, gotConn bool) failureKind {
	if gotConn {
		return responseLost
	}
	if isNetworkError(err) || isTLSError(err) || errors.Is(err, context.Canceled) {
		return neverSent
	}
	return notRetryable
}

// isResponseLost reports whether err means that the complete response was not received: the connection broke or timed
// out while the body was being read. A complete body that cannot be parsed is not one.
func isResponseLost(resp *req.Response, err error) bool {
	return resp.Response == nil || isNetworkError(err)
}

// isNetworkError reports whether err is a failure of the connection: a dial, read or write error, a DNS failure, a
// timeout, or a connection closed too early. *url.Error is unwrapped first: it implements net.Error for every error,
// an unsupported protocol scheme included.
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
