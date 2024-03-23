package responses

type ReportSettingsResponse struct {
	Number  string `json:"number"`
	Name    string `json:"name"`
	Columns int64  `json:"columns"`
	Status  string `json:"status"`
}
