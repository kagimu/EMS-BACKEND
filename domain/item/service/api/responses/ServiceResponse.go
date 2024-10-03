package responses

type ItemResponse struct {
	Number          string  `json:"number"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Type            string  `json:"type"`
	Cost            float64 `json:"cost"`
	SerialCode      string  `json:"serial_code"`
	Prefix          string  `json:"prefix"`
	PaymentMode     string  `json:"payment_mode"`
	Mandatory       bool    `json:"mandatory"`
	SchoolNumber    string  `json:"school_number"`
	UserNumber      string  `json:"user_number"`
	RepeatCycle     string  `json:"repeat_cycle"`
	OneTimePurchase bool    `json:"one_time_purchase"`
}
type ServiceResponse struct {
	ItemNumber string                  `json:"item_number"`
	Item       ItemResponse            `json:"item"`
	Campuses   []CampusServiceResponse `json:"campus_services"`
}
type CampusServiceResponse struct {
	Number       string `json:"number"`
	CampusNumber string `json:"campus_number"`
	Name         string `json:"name"`
}
