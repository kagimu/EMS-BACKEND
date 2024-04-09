package responses

type FeedingResponse struct {
	Number               string `json:"number"`
	FeedingSessionNumber string `json:"feeding_session_number"`
	FeedingSessionName   string `json:"feeding_session_name"`
	FeedingSessionStatus string `json:"feeding_session_status"`
	FeedingSessionCode   uint64 `json:"feeding_session_code"`
	Date                 string `json:"date"`
	StaffNumber          string `json:"staff_number"`
	StaffName            string `json:"staff_name"`
	StockingNumber       string `json:"stocking_number"`
	StockingQuantity     int    `json:"stocking_quantity"`
}
