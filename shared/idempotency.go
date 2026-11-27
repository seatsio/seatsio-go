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

type IdempotencyDetails struct {
	IdempotencyKey         string
	IdempotencyKeyHonored  *bool
	Replayed               bool
	outcomeUnknown         bool
	safeToRetryWithSameKey bool
}

func (d *IdempotencyDetails) IsOutcomeUnknown() bool {
	return d.outcomeUnknown
}

func (d *IdempotencyDetails) IsSafeToRetryWithSameKey() bool {
	return d.safeToRetryWithSameKey
}

func (d *IdempotencyDetails) idempotencyDetails() *IdempotencyDetails {
	return d
}

type idempotencyFailure interface {
	idempotencyDetails() *IdempotencyDetails
}

func idempotencyDetailsOf(err error) (*IdempotencyDetails, bool) {
	var failure idempotencyFailure
	if errors.As(err, &failure) {
		return failure.idempotencyDetails(), true
	}
	return nil, false
}

type RequestError struct {
	IdempotencyDetails
	Err error
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

func IsOutcomeUnknown(err error) bool {
	details, ok := idempotencyDetailsOf(err)
	return ok && details.outcomeUnknown
}

func IsSafeToRetryWithSameKey(err error) bool {
	details, ok := idempotencyDetailsOf(err)
	return ok && details.safeToRetryWithSameKey
}

func IdempotencyKeyOf(err error) (IdempotencyKey, bool) {
	details, ok := idempotencyDetailsOf(err)
	if !ok {
		return IdempotencyKey{}, false
	}
	key, parseErr := ParseIdempotencyKey(details.IdempotencyKey)
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
