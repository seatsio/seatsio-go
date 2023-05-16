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
}

func (m *SeatsioError) Error() string {
	return m.Message
}
