package responses

type DayItemsResponse struct {
	Day    string   `json:"day"`
	Active bool     `json:"active"`
	Error  string   `json:"error"`
	Items  []string `json:"items"`
	Name   string   `json:"name"`
}
