package responses

type AggregateScaleResponse struct {
	Code         int64   `json:"code"`
	Number       string  `json:"number"`
	MinValue     float64 `json:"min_value"`
	MaxValue     float64 `json:"max_value"`
	Name         string  `json:"name"`
	Weight       float64 `json:"weight"`
	Illustration string  `json:"illustration"`
	Description  string  `json:"description"`
	ShortName    string  `json:"short_name"`
	ScaleName    string  `json:"scale_name"`
}

type ScaleResponse struct {
	ScaleNumber string `json:"scale_number"`
	ScaleName   string `json:"scale_name"`
	Count       int64  `json:"count"`
}

type DetailedScaleResponse struct {
	ScaleResponse
	Grades []AggregateScaleResponse `json:"grades"`
}
