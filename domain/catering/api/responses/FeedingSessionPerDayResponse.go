package responses

type FeedingSessionPerDayResponse struct {
	Number               string `json:"number"`
	FeedingSessionNumber string `json:"feeding_session_number"`
	Day                  string `json:"day"`
	FeedingSessionName   string `json:"feeding_session_name"`
}
