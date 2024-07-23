package responses

type FeedingSessionResponse struct {
	Name         string `json:"name"`
	Code         uint64 `json:"code"`
	Number       string `json:"number"`
	AddedBy      string `json:"added_by"`
	SchoolNumber string `json:"school_number"`
}
