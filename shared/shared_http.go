package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/imroc/req/v3"
)

type AdditionalHeader func(headers *map[string]string)

func ApiClient(secretKey string, baseUrl string, additionalHeaders ...AdditionalHeader) *req.Client {
	client := req.C().SetBaseURL(baseUrl).
		SetCommonBasicAuth(secretKey, "").
		SetTimeout(10 * time.Second).
		SetCommonRetryCount(5).
		SetCommonRetryInterval(interruptibleRetryInterval(backoffInterval(400*time.Millisecond, 10*time.Second))).
		SetCommonRetryCondition(retryCondition(false)).
		WrapRoundTripFunc(addIdempotencyKey).
		SetHTTP2ReadIdleTimeout(30 * time.Second)
	headers := make(map[string]string)
	for _, opt := range additionalHeaders {
		opt(&headers)
	}
	for key, value := range headers {
		client.SetCommonHeader(key, value)
	}
	return client
}

// SetRetryWritesOnConnectionFailure controls whether POST, PUT, PATCH and DELETE requests are retried (with the same
// idempotency key) when the connection fails or times out before the complete response was received. GET requests are
// always retried in that case. 429 and 409 IDEMPOTENCY_KEY_IN_USE responses are retried regardless of this setting.
func SetRetryWritesOnConnectionFailure(client *req.Client, enabled bool) {
	client.SetCommonRetryCondition(retryCondition(enabled))
}

func retryCondition(retryWritesOnConnectionFailure bool) req.RetryConditionFunc {
	return func(resp *req.Response, err error) bool {
		if err != nil {
			if resp.Request != nil && resp.Request.Context().Err() != nil {
				return false
			}
			return isConnectionFailure(resp, err) && (retryWritesOnConnectionFailure || isSafeMethod(resp.Request.Method))
		}
		return resp.StatusCode == 429 || isKeyInUse(resp)
	}
}

// isConnectionFailure reports whether err means the complete response was not received: either no response at all, or
// the connection broke or timed out while the body was being read. A complete body that cannot be parsed is not one.
func isConnectionFailure(resp *req.Response, err error) bool {
	if resp.Response == nil {
		return true
	}
	var netErr net.Error
	return errors.Is(err, io.ErrUnexpectedEOF) || errors.As(err, &netErr)
}

// backoffInterval is req's capped exponential backoff with jitter (SetCommonRetryBackoffInterval), which req does not
// export.
func backoffInterval(min, max time.Duration) req.GetRetryIntervalFunc {
	return func(_ *req.Response, attempt int) time.Duration {
		capped := math.Min(float64(max), float64(min)*math.Exp2(float64(attempt)))
		half := int64(capped / 2)
		return time.Duration(half + rand.Int63n(half))
	}
}

// interruptibleRetryInterval waits for the interval itself and then tells req not to wait any longer. req waits between
// retries with time.Sleep, which a cancelled or expired context cannot interrupt, so a 429 or 409 retry could otherwise
// keep the caller waiting up to 10 seconds after the context ended. When the context ends, the next attempt fails
// before anything is sent, and the retry condition stops.
func interruptibleRetryInterval(interval req.GetRetryIntervalFunc) req.GetRetryIntervalFunc {
	return func(resp *req.Response, attempt int) time.Duration {
		wait := interval(resp, attempt)
		if resp == nil || resp.Request == nil {
			return wait
		}
		timer := time.NewTimer(wait)
		defer timer.Stop()
		select {
		case <-timer.C:
		case <-resp.Request.Context().Done():
		}
		return 0
	}
}

func isKeyInUse(resp *req.Response) bool {
	if resp.StatusCode != http.StatusConflict {
		return false
	}
	apiError, err := firstApiErrorIfAny(resp)
	return err == nil && apiError != nil && apiError.Code == "IDEMPOTENCY_KEY_IN_USE"
}

func firstApiErrorIfAny(resp *req.Response) (*SeatsioErrorTO, error) {
	errorResponse := &SeatsioErrorResponse{}
	if err := json.Unmarshal(resp.Bytes(), errorResponse); err != nil {
		return nil, err
	}
	if len(errorResponse.Errors) == 0 {
		return nil, nil
	}
	return &errorResponse.Errors[0], nil
}

func unexpectedResponseError(resp *req.Response) error {
	return fmt.Errorf("server returned error %v. Body: %v", resp.StatusCode, string(resp.Bytes()))
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

// idempotencyKeyWireHeader is the header the key is sent in. It is deliberately not in canonical form: the transport
// treats a request with a canonical Idempotency-Key header as replayable, and silently resends it on a reused
// connection that breaks before the response arrives, even when retries are disabled. That would execute the write a
// second time for a company without idempotency. Header names are case-insensitive, so the server reads it all the same.
const idempotencyKeyWireHeader = "idempotency-key"

// addIdempotencyKey runs right before each attempt is handed to the transport, after the client's common headers were
// merged into the request, so it also sees a key set as a common header.
func addIdempotencyKey(rt req.RoundTripper) req.RoundTripFunc {
	return func(r *req.Request) (*req.Response, error) {
		if !isSafeMethod(r.Method) {
			setIdempotencyKey(r)
		}
		return rt.RoundTrip(r)
	}
}

func setIdempotencyKey(r *req.Request) {
	if r.Headers == nil {
		r.Headers = make(http.Header)
	}
	if callerKey, ok := r.Headers[IdempotencyKeyHeader]; ok {
		delete(r.Headers, IdempotencyKeyHeader)
		r.Headers[idempotencyKeyWireHeader] = callerKey
	}
	if len(r.Headers[idempotencyKeyWireHeader]) > 0 {
		return
	}
	key, ok := IdempotencyKeyFromContext(r.Context())
	if !ok {
		key = NewIdempotencyKey()
	}
	r.Headers[idempotencyKeyWireHeader] = []string{key.String()}
}

func AssertOk[T interface{}](result *req.Response, err error, data *T) (*T, error) {
	err = AssertOkWithoutResult(result, err)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func AssertOkNoBody(result *req.Response, err error) error {
	err = AssertOkWithoutResult(result, err)
	if err != nil {
		return err
	}
	return nil
}

func AssertOkMap[T interface{}](result *req.Response, err error, data map[string]T) (map[string]T, error) {
	err = AssertOkWithoutResult(result, err)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func AssertOkArray[T interface{}](result *req.Response, err error, data *[]T) ([]T, error) {
	err = AssertOkWithoutResult(result, err)
	if err != nil {
		return nil, err
	}
	return *data, nil
}

func AssertOkWithoutResult(result *req.Response, err error) error {
	key := sentIdempotencyKey(result)
	if err != nil {
		return requestError(key, isOutcomeUnknown(result, err), err)
	}
	if !result.IsSuccessState() {
		outcomeUnknown := result.StatusCode >= 500
		if !strings.Contains(result.GetHeader("content-type"), "application/json") {
			return requestError(key, outcomeUnknown, unexpectedResponseError(result))
		}
		apiError, err := firstApiErrorIfAny(result)
		if err != nil {
			return requestError(key, outcomeUnknown, err)
		}
		if apiError == nil {
			return requestError(key, outcomeUnknown, unexpectedResponseError(result))
		}
		return &SeatsioError{Message: apiError.Message, Code: apiError.Code, IdempotencyKey: key, outcomeUnknown: outcomeUnknown && key != ""}
	}
	return nil
}

// isOutcomeUnknown reports whether a failed request may or may not have been executed by the server: no status was
// received, the status was a 5xx, or the status was a 2xx but the connection broke or timed out while the body was
// being read.
func isOutcomeUnknown(result *req.Response, err error) bool {
	if result == nil || result.Response == nil {
		return true
	}
	if result.StatusCode >= 500 {
		return true
	}
	return result.IsSuccessState() && isConnectionFailure(result, err)
}

func requestError(key string, outcomeUnknown bool, err error) error {
	if key == "" {
		return err
	}
	return &RequestError{IdempotencyKey: key, Err: err, outcomeUnknown: outcomeUnknown}
}

func WithAdditionalHeader(key string, value string) AdditionalHeader {
	return func(headers *map[string]string) {
		(*headers)[key] = value
	}
}
