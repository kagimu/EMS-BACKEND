package requests

type CalendarEventRequest struct {
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	CalendarNumber string `json:"calendar_number"`
}

type GeneralCalendarRequest struct {
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	EntityNumber string `json:"entity_number"`
}

type EntityRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}
