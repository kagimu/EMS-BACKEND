package response

type IndicatorResponse struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Code        string `json:"code"`
	InputType   string `json:"input_type"`
	Description string `json:"description"`
	Number      string `json:"number"`
}
