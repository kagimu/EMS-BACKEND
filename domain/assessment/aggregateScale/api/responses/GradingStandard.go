package responses

type GradingStandardResponse struct {
	Code        int64  `json:"code" `
	Number      string `json:"number" `
	Name        string `json:"name"`
	Description string `json:"description" `
	Type        string `json:"type"`
}
