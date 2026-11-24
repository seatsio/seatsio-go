package shared

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/imroc/req/v3"
)

type AdditionalHeader func(headers *map[string]string)

func ApiClient(secretKey string, baseUrl string, additionalHeaders ...AdditionalHeader) *req.Client {
	return ApiClientWithRetries(secretKey, baseUrl, NewRetries(SystemClock), additionalHeaders...)
}

// ApiClientWithRetries creates an API client that keeps its retry state in retries. Don't share retries between clients
// with different secret keys.
func ApiClientWithRetries(secretKey string, baseUrl string, retries *Retries, additionalHeaders ...AdditionalHeader) *req.Client {
	client := req.C().SetBaseURL(baseUrl).
		SetCommonBasicAuth(secretKey, "").
		SetTimeout(10 * time.Second).
		SetCommonRetryCount(int(retries.maxRetries.Load())).
		SetCommonRetryInterval(retries.interval).
		SetCommonRetryCondition(retryCondition).
		WrapRoundTripFunc(retries.sendWrites).
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

// retryCondition is req's retry condition. req retries GET, HEAD, OPTIONS and TRACE requests: on a 429, and on a
// connection failure while the context is active. Other requests are retried by Retries.sendWrites, not by req.
func retryCondition(resp *req.Response, err error) bool {
	if resp.Request == nil || !isSafeMethod(resp.Request.Method) {
		return false
	}
	if err != nil {
		return resp.Request.Context().Err() == nil && isResponseLost(resp, err)
	}
	return resp.StatusCode == 429
}

// backoff is req's capped exponential backoff with jitter (SetCommonRetryBackoffInterval), which req does not export.
// attempt starts at 1.
func backoff(min, max time.Duration) func(attempt int) time.Duration {
	return func(attempt int) time.Duration {
		capped := math.Min(float64(max), float64(min)*math.Exp2(float64(attempt)))
		half := int64(capped / 2)
		return time.Duration(half + rand.Int63n(half))
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

// setIdempotencyKey runs right before the first attempt is handed to the transport, after the client's common headers
// were merged into the request, so it also sees a key set as a common header. Retries reuse the request, and so the key.
func setIdempotencyKey(r *req.Request) error {
	if r.Headers == nil {
		r.Headers = make(http.Header)
	}
	if callerKey, ok := r.Headers[IdempotencyKeyHeader]; ok {
		delete(r.Headers, IdempotencyKeyHeader)
		r.Headers[idempotencyKeyWireHeader] = callerKey
	}
	if len(r.Headers[idempotencyKeyWireHeader]) > 0 {
		return nil
	}
	key, ok := IdempotencyKeyFromContext(r.Context())
	if ok && key.IsZero() {
		return errZeroIdempotencyKey
	}
	if !ok {
		key = NewIdempotencyKey()
	}
	r.Headers[idempotencyKeyWireHeader] = []string{key.String()}
	return nil
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
		if key == "" {
			return err
		}
		requestErr := &RequestError{IdempotencyKey: key, Err: err}
		if result != nil && result.Response != nil {
			status := idempotencyKeyStatusOf(result.Header)
			requestErr.IdempotencyKeyHonored, requestErr.Replayed = status.keyHonoredOrNil(), status == statusReplayed
		}
		if w := writeOf(result); w != nil && (result.Response == nil || isResponseLost(result, err)) {
			requestErr.outcomeUnknown, requestErr.safeToRetryWithSameKey = w.failureWithoutResponse()
		}
		return requestErr
	}
	if result.IsSuccessState() {
		return nil
	}
	apiError, parseErr := firstApiError(result)
	if key == "" {
		if parseErr != nil {
			return parseErr
		}
		return &SeatsioError{Message: apiError.Message, Code: apiError.Code}
	}
	status := idempotencyKeyStatusOf(result.Header)
	w := writeOf(result)
	outcomeUnknown := result.StatusCode >= 500 || isKeyInUse(result)
	keyHonored := status.keyHonored()
	if !status.known() {
		keyHonored = w != nil && w.routeHonorsKeys()
	}
	safeToRetryWithSameKey := outcomeUnknown && keyHonored && status != statusReplayed
	if result.StatusCode == http.StatusTooManyRequests && w != nil && w.mayHaveRun {
		// The 429 says nothing about an earlier attempt of this call, which may have run.
		outcomeUnknown, safeToRetryWithSameKey = true, w.routeHonorsKeys() && w.startsWithinRetryWindow(0)
	}
	if parseErr != nil {
		return &RequestError{
			IdempotencyKey:         key,
			IdempotencyKeyHonored:  status.keyHonoredOrNil(),
			Replayed:               status == statusReplayed,
			Err:                    parseErr,
			outcomeUnknown:         outcomeUnknown,
			safeToRetryWithSameKey: safeToRetryWithSameKey,
		}
	}
	return &SeatsioError{
		Message:                apiError.Message,
		Code:                   apiError.Code,
		IdempotencyKey:         key,
		IdempotencyKeyHonored:  status.keyHonoredOrNil(),
		Replayed:               status == statusReplayed,
		outcomeUnknown:         outcomeUnknown,
		safeToRetryWithSameKey: safeToRetryWithSameKey,
	}
}

// firstApiError returns the first error of a seats.io error response, or an error that describes the response when it
// isn't one: not JSON, unparsable, or without errors. It never panics.
func firstApiError(result *req.Response) (*SeatsioErrorTO, error) {
	if !strings.Contains(result.GetHeader("content-type"), "application/json") {
		return nil, unexpectedResponseError(result)
	}
	apiError, err := firstApiErrorIfAny(result)
	if err != nil {
		return nil, err
	}
	if apiError == nil {
		return nil, unexpectedResponseError(result)
	}
	return apiError, nil
}

func WithAdditionalHeader(key string, value string) AdditionalHeader {
	return func(headers *map[string]string) {
		(*headers)[key] = value
	}
}
