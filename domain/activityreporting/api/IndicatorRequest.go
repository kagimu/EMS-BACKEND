package api

type IndicatorRequest struct {
	Name         string `json:"name"`
	EntityNumber string `json:"entityNumber"`
	Description  string `json:"description"`
	Status       string `json:"status"`
}
