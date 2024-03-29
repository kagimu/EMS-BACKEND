package responses

type ReportSettingsResponse struct {
	Number  string `json:"number"`
	Name    string `json:"name"`
	Columns int64  `json:"columns"`
	Status  string `json:"status"`
	Code    int64  `json:"code"`
}

type ReportColumnResponse struct {
	Number          string  `json:"number"`
	Name            string  `json:"name"`
	PercentageShare float64 `json:"percentage_share"`
	Code            int64   `json:"code"`
}

type ReportColumnAssessmentResponse struct {
	Number           string `json:"number"`
	AssessmentArea   string `json:"assessment_area"`
	AssessmentMetric string `json:"assessment_metric"`
	Code             int64  `json:"code"`
	ReportColumn     string `json:"report_column"`
}

type SummarizedFinalReportResponse struct {
	Student   string `json:"student"`
	Grade     string `json:"grade"`
	Aggregate string `json:"aggregate"`
	Class     string `json:"class"`
}
