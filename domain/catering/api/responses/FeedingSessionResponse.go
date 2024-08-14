package responses

type FeedingSessionResponse struct {
	FeedingSessionNumber string     `json:"feeding_session_number"`
	Name                 string     `json:"name"`
	AddedBy              string     `json:"added_by"`
	SchoolNumber         string     `json:"school_number"`
	PackCost             float64    `json:"pack_cost"`
	PackItems            []PackItem `json:"pack_items"`
}

type PackItem struct {
	Quantity          float64               `json:"consumption_quantity"`
	Name              string                `json:"name"`
	Type              string                `json:"type"`
	ItemNumber        string                `json:"item_number"`
	UnitOfConsumption UnitOfMeasureResponse `json:"unit_of_consumption"`
}
