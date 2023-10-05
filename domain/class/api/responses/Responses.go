package responses

type (
	ClassLevelInfo struct {
		Number       string       `json:"number"`
		ClassNumber  string       `json:"class_number"`
		Name         string       `json:"name"`
		Abbreviation string       `json:"abbreviation"`
		Streams      []StreamInfo `json:"streams"`
	}
	StreamInfo struct {
		StreamNumber string `json:"stream_number"`
		Name         string `json:"name"`
	}
)
