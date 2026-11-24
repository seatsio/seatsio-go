package shared

type SeatsioErrorTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SeatsioErrorResponse struct {
	Errors []SeatsioErrorTO `json:"errors"`
}

type SeatsioError struct {
	Code    string
	Message string
	// IdempotencyKey is the key a POST, PUT, PATCH or DELETE request was sent with; empty for other requests.
	IdempotencyKey string
	// IdempotencyKeyHonored says whether the server honored the key, according to the Idempotency-Key-Status header of
	// the response; nil when the response didn't say.
	IdempotencyKeyHonored *bool
	// Replayed is true when the server didn't execute the request again, but replayed the recorded response of an
	// earlier request with the same key.
	Replayed               bool
	outcomeUnknown         bool
	safeToRetryWithSameKey bool
}

func (m *SeatsioError) Error() string {
	return m.Message
}

// IsOutcomeUnknown reports whether the request may or may not have been executed. See shared.IsOutcomeUnknown.
func (m *SeatsioError) IsOutcomeUnknown() bool {
	return m.outcomeUnknown
}

// IsSafeToRetryWithSameKey reports whether the same call, made again with IdempotencyKey, can't execute the request
// twice and may give a different result. See shared.IsSafeToRetryWithSameKey.
func (m *SeatsioError) IsSafeToRetryWithSameKey() bool {
	return m.safeToRetryWithSameKey
}
