package charts

type ChartValidationResult struct {
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}
