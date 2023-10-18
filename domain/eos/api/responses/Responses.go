package responses

import "time"

type Calendar struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Number       string    `json:"number"`
	EntityNumber string    `json:"entity_number"`
}
