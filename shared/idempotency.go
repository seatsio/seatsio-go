package shared

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/imroc/req/v3"
)

const IdempotencyKeyHeader = "Idempotency-Key"

// IdempotencyKeyStatusHeader is the response header in which the server says what it did with the Idempotency-Key of a
// POST or DELETE request that reached a route.
const IdempotencyKeyStatusHeader = "Idempotency-Key-Status"

type idempotencyKeyContextKey struct{}

type IdempotencyKey struct{ id uuid.UUID }

// v4UUID is deliberately strict: uuid.Parse also accepts braces, a urn:uuid: prefix and the form without dashes.
var v4UUID = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func NewIdempotencyKey() IdempotencyKey {
	return IdempotencyKey{id: uuid.New()}
}

// ParseIdempotencyKey accepts a v4 UUID in its canonical form (8-4-4-4-12 hexadecimal digits, any case) and rejects
// anything else.
func ParseIdempotencyKey(s string) (IdempotencyKey, error) {
	if !v4UUID.MatchString(s) {
		return IdempotencyKey{}, fmt.Errorf("idempotency key %q is not a v4 UUID", s)
	}
	return IdempotencyKey{id: uuid.MustParse(s)}, nil
}

func (k IdempotencyKey) String() string {
	return k.id.String()
}

func (k IdempotencyKey) IsZero() bool {
	return k.id == uuid.Nil
}

// WithIdempotencyKey returns a context that makes every POST, PUT, PATCH and DELETE request made with it send key. Use
// it for one operation only. A zero key makes the request fail before it is sent.
func WithIdempotencyKey(ctx context.Context, key IdempotencyKey) context.Context {
	return context.WithValue(ctx, idempotencyKeyContextKey{}, key)
}

// IdempotencyKeyFromContext returns the key that WithIdempotencyKey put on ctx, if any. The key may be zero.
func IdempotencyKeyFromContext(ctx context.Context) (IdempotencyKey, bool) {
	key, ok := ctx.Value(idempotencyKeyContextKey{}).(IdempotencyKey)
	return key, ok
}

var errZeroIdempotencyKey = errors.New("the idempotency key in the context is the zero value; create one with NewIdempotencyKey or ParseIdempotencyKey")

// idempotencyKeyStatus is the value of the Idempotency-Key-Status response header. The zero value means the header is
// absent or has a value this client doesn't know; both say nothing about the route.
type idempotencyKeyStatus int

const (
	statusUnknown idempotencyKeyStatus = iota
	statusExecuted
	statusReplayed
	statusInUse
	statusReused
	statusIgnored
)

// idempotencyKeyStatusOf reads the Idempotency-Key-Status header. The four values that mean the route honors keys are
// listed explicitly, so a future value or a proxy that strips the header errs towards not retrying.
func idempotencyKeyStatusOf(header http.Header) idempotencyKeyStatus {
	switch strings.ToLower(strings.TrimSpace(header.Get(IdempotencyKeyStatusHeader))) {
	case "executed":
		return statusExecuted
	case "replayed":
		return statusReplayed
	case "in-use":
		return statusInUse
	case "reused":
		return statusReused
	case "ignored":
		return statusIgnored
	default:
		return statusUnknown
	}
}

func (s idempotencyKeyStatus) known() bool {
	return s != statusUnknown
}

// keyHonored reports whether the route honors idempotency keys. Only meaningful for a known status.
func (s idempotencyKeyStatus) keyHonored() bool {
	return s != statusUnknown && s != statusIgnored
}

// keyHonoredOrNil is keyHonored for a known status, and nil otherwise.
func (s idempotencyKeyStatus) keyHonoredOrNil() *bool {
	if !s.known() {
		return nil
	}
	honored := s.keyHonored()
	return &honored
}

// RequestError is returned when a POST, PUT, PATCH or DELETE request fails without an error response from the API:
// a connection failure, a timeout, or a response that is not a seats.io error. It wraps the original error.
type RequestError struct {
	// IdempotencyKey is the key the request was sent with.
	IdempotencyKey string
	// IdempotencyKeyHonored says whether the server honored the key, according to the Idempotency-Key-Status header of
	// the response; nil when there was no response or the response didn't say.
	IdempotencyKeyHonored *bool
	// Replayed is true when the server didn't execute the request again, but replayed the recorded response of an
	// earlier request with the same key.
	Replayed               bool
	Err                    error
	outcomeUnknown         bool
	safeToRetryWithSameKey bool
}

func (e *RequestError) Error() string {
	return e.Err.Error()
}

func (e *RequestError) Unwrap() error {
	return e.Err
}

// Timeout forwards to the wrapped error, so os.IsTimeout and net.Error checks keep working.
func (e *RequestError) Timeout() bool {
	var t interface{ Timeout() bool }
	return errors.As(e.Err, &t) && t.Timeout()
}

// Temporary forwards to the wrapped error, so net.Error checks keep working.
func (e *RequestError) Temporary() bool {
	var t interface{ Temporary() bool }
	return errors.As(e.Err, &t) && t.Temporary()
}

// IsOutcomeUnknown reports whether the request may or may not have been executed. See shared.IsOutcomeUnknown.
func (e *RequestError) IsOutcomeUnknown() bool {
	return e.outcomeUnknown
}

// IsSafeToRetryWithSameKey reports whether the same call, made again with IdempotencyKey, can't execute the request
// twice and may give a different result. See shared.IsSafeToRetryWithSameKey.
func (e *RequestError) IsSafeToRetryWithSameKey() bool {
	return e.safeToRetryWithSameKey
}

// IsOutcomeUnknown reports whether err is a failed POST, PUT, PATCH or DELETE request that the server may or may not
// have executed: the response was lost after the request was sent (a timeout, a connection reset, an incomplete body),
// the server answered with a 5xx, or it answered 409 IDEMPOTENCY_KEY_IN_USE after the client ran out of retries. It is
// false when the request never reached the server, for every other 4xx, and for errors of GET requests.
func IsOutcomeUnknown(err error) bool {
	var requestErr *RequestError
	if errors.As(err, &requestErr) {
		return requestErr.outcomeUnknown
	}
	var apiErr *SeatsioError
	if errors.As(err, &apiErr) {
		return apiErr.outcomeUnknown
	}
	return false
}

// IsSafeToRetryWithSameKey reports whether the call that returned err, made again with the idempotency key of err
// (see IdempotencyKeyOf and WithIdempotencyKey), can't execute the request twice and may give a different result. That
// holds when the request never reached the server, or when its outcome is unknown and the server honors idempotency
// keys for it: known from the response or, when the response was lost, from the last response of the same operation.
// Retry promptly: the server keeps a key for 10 minutes.
func IsSafeToRetryWithSameKey(err error) bool {
	var requestErr *RequestError
	if errors.As(err, &requestErr) {
		return requestErr.safeToRetryWithSameKey
	}
	var apiErr *SeatsioError
	if errors.As(err, &apiErr) {
		return apiErr.safeToRetryWithSameKey
	}
	return false
}

func IdempotencyKeyOf(err error) (IdempotencyKey, bool) {
	var raw string
	var requestErr *RequestError
	var apiErr *SeatsioError
	if errors.As(err, &requestErr) {
		raw = requestErr.IdempotencyKey
	} else if errors.As(err, &apiErr) {
		raw = apiErr.IdempotencyKey
	}
	key, parseErr := ParseIdempotencyKey(raw)
	return key, parseErr == nil
}

func sentIdempotencyKey(result *req.Response) string {
	if result == nil || result.Request == nil {
		return ""
	}
	if values := result.Request.Headers[idempotencyKeyWireHeader]; len(values) > 0 {
		return values[0]
	}
	return ""
}
