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

const IdempotencyKeyStatusHeader = "Idempotency-Key-Status"

type idempotencyKeyContextKey struct{}

type IdempotencyKey struct{ id uuid.UUID }

var v4UUID = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func NewIdempotencyKey() IdempotencyKey {
	return IdempotencyKey{id: uuid.New()}
}

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

func WithIdempotencyKey(ctx context.Context, key IdempotencyKey) context.Context {
	return context.WithValue(ctx, idempotencyKeyContextKey{}, key)
}

func IdempotencyKeyFromContext(ctx context.Context) (IdempotencyKey, bool) {
	key, ok := ctx.Value(idempotencyKeyContextKey{}).(IdempotencyKey)
	return key, ok
}

var errZeroIdempotencyKey = errors.New("the idempotency key in the context is the zero value; create one with NewIdempotencyKey or ParseIdempotencyKey")

type idempotencyKeyStatus int

const (
	statusUnknown idempotencyKeyStatus = iota
	statusExecuted
	statusReplayed
	statusInUse
	statusReused
	statusIgnored
)

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

func (s idempotencyKeyStatus) keyHonored() bool {
	return s != statusUnknown && s != statusIgnored
}

func (s idempotencyKeyStatus) keyHonoredOrNil() *bool {
	if !s.known() {
		return nil
	}
	honored := s.keyHonored()
	return &honored
}

type RequestError struct {
	IdempotencyKey         string
	IdempotencyKeyHonored  *bool
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

func (e *RequestError) Timeout() bool {
	var t interface{ Timeout() bool }
	return errors.As(e.Err, &t) && t.Timeout()
}

func (e *RequestError) Temporary() bool {
	var t interface{ Temporary() bool }
	return errors.As(e.Err, &t) && t.Temporary()
}

func (e *RequestError) IsOutcomeUnknown() bool {
	return e.outcomeUnknown
}

func (e *RequestError) IsSafeToRetryWithSameKey() bool {
	return e.safeToRetryWithSameKey
}

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
