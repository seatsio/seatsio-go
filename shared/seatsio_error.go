package shared

type SeatsioErrorTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SeatsioErrorResponse struct {
	Errors []SeatsioErrorTO `json:"errors"`
}

type SeatsioError struct {
	Code                   string
	Message                string
	IdempotencyKey         string
	IdempotencyKeyHonored  *bool
	Replayed               bool
	outcomeUnknown         bool
	safeToRetryWithSameKey bool
}

func (m *SeatsioError) Error() string {
	return m.Message
}

func (m *SeatsioError) IsOutcomeUnknown() bool {
	return m.outcomeUnknown
}

func (m *SeatsioError) IsSafeToRetryWithSameKey() bool {
	return m.safeToRetryWithSameKey
}
