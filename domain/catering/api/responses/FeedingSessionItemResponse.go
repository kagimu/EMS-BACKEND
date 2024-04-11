package responses

type FeedingSessionItemResponse struct {
	Number               string `json:"number"`
	FeedingSessionNumber string `json:"feeding_session_number"`
	FeedingSessionName   string `json:"feeding_session_name"`
	FeedingSessionStatus string `json:"feeding_session_status"`
	MenuItemNumber       string `json:"menu_item_number"`
	MenuItemName         string `json:"menu_item_name"`
}
