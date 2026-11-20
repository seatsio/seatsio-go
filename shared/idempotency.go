package shared

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/imroc/req/v3"
)

const IdempotencyKeyHeader = "Idempotency-Key"

type idempotencyKeyContextKey struct{}

type IdempotencyKey struct{ id uuid.UUID }

func NewIdempotencyKey() IdempotencyKey {
	return IdempotencyKey{id: uuid.New()}
}

func ParseIdempotencyKey(s string) (IdempotencyKey, error) {
	id, err := uuid.Parse(s)
	if err != nil || id.Version() != 4 || id.Variant() != uuid.RFC4122 {
		return IdempotencyKey{}, fmt.Errorf("idempotency key %q is not a v4 UUID", s)
	}
	return IdempotencyKey{id: id}, nil
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
	return key, ok && !key.IsZero()
}

// RequestError is returned when a POST, PUT, PATCH or DELETE request fails without an error response from the API:
// a connection failure, a timeout, or a response that is not a seats.io error. It wraps the original error.
type RequestError struct {
	IdempotencyKey string
	Err            error
	outcomeUnknown bool
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

// IsOutcomeUnknown reports whether err is a failed POST, PUT, PATCH or DELETE request that the server may or may not
// have executed: the connection failed or timed out before a status was received or while a successful response was
// being read, or the server answered with a 5xx.
// Only such a request is worth retrying with the same idempotency key; any other error happens again on a retry.
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
