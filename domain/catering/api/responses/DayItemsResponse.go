package responses

type DayItemsResponse struct {
	Day   string   `json:"day"`
	Items []string `json:"items"`
}
