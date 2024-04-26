package responses

type FeedingGroupResponse struct {
	Number       string `json:"number"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	CampusNumber string `json:"campus_number"`
	CampusName   string `json:"campus_name"`
	UserName     string `json:"user_name"`
	Photo        string `json:"photo"`
}

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

type FeedingSessionItemResponse struct {
	Number               string `json:"number"`
	FeedingSessionNumber string `json:"feeding_session_number"`
	FeedingSessionName   string `json:"feeding_session_name"`
	FeedingSessionStatus string `json:"feeding_session_status"`
	MenuItemNumber       string `json:"menu_item_number"`
	MenuItemName         string `json:"menu_item_name"`
}

type FeedingSessionPerDayResponse struct {
	Number               string `json:"number"`
	FeedingSessionNumber string `json:"feeding_session_number"`
	Day                  string `json:"day"`
	FeedingSessionName   string `json:"feeding_session_name"`
}

type FeedingUserGroupResponse struct {
	UserNumber         string `json:"user_number"`
	Number             string `json:"number"`
	Gender             string `json:"gender"`
	Photo              string `json:"photo"`
	Name               string `json:"name"`
	FeedingGroupNumber string `json:"feeding_group_number"`
	FeedingGroupName   string `json:"feeding_group_name"`
	Status             string `json:"status"`
}

type ItemAssemblyResponse struct {
	Number              string  `json:"number"`
	PrimaryItemNumber   string  `json:"primary_item_number"`
	SecondaryItemNumber string  `json:"secondary_item_number"`
	Quantity            float64 `json:"quantity"`
	PrimaryItemName     string  `json:"primary_item_name"`
	SecondaryItemName   string  `json:"secondary_item_name"`
}

type MenuItemResponse struct {
	Number              string `json:"number"`
	Name                string `json:"name"`
	Code                uint64 `json:"code"`
	Type                string `json:"type"`
	CampusNumber        string `json:"campus_number"`
	CampusName          string `json:"campus_name"`
	Phone               string `json:"phone"`
	Email               string `json:"email"`
	Location            string `json:"location"`
	UnitOfMeasureNumber string `json:"unit_of_measure_number"`
	UnitOfMeasureName   string `json:"unit_of_measure_name"`
	UnitOfMeasureCode   uint64 `json:"unit_of_measure_code"`
	Parts               []Part `json:"parts"`
}

type Part struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Code     uint64  `json:"code"`
	Number   string  `json:"number"`
}

type StockingResponse struct {
	Number         string  `json:"name"`
	Quantity       int     `json:"quantity"`
	UnitCost       float64 `json:"unit_cost"`
	MenuItemNumber string  `json:"menu_item_number"`
	MenuItemName   string  `json:"menu_item_name"`
	CampusNumber   string  `json:"campus_number"`
	CampusName     string  `json:"campus_name"`
	MenuItemType   string  `json:"menu_item_type"`
}

type FeedingGroupSessionResponse struct {
	UserNumber           string `json:"user_number"`
	Number               string `json:"number"`
	FeedingGroupNumber   string `json:"feeding_group_number"`
	FeedingGroupName     string `json:"feeding_group_name"`
	Status               string `json:"status"`
	FeedingSessionNumber string `json:"feeding_session_number"`
	FeedingSessionName   string `json:"feeding_session_name"`
	FeedingSessionStatus string `json:"feeding_session_status"`
	FeedingSessionCampus string `json:"feeding_session_campus"`
}

type FeedingsPerGroupResponse struct {
	Number               string `json:"number"`
	GroupName            string `json:"group_name"`
	GroupStatus          string `json:"group_status"`
	Day                  string `json:"day"`
	FeedingSessionNumber string `json:"feeding_session_number"`
	FeedingSessionName   string `json:"feeding_session_name"`
	FeedingSessionStatus string `json:"feeding_session_status"`
	FeedingSessionCode   uint64 `json:"feeding_session_code"`
}

type FeedingGroupPerDayResponse struct {
	Number             string `json:"number"`
	FeedingGroupNumber string `json:"feeding_group_number"`
	FeedingGroupName   string `json:"feeding_group_name"`
	Day                string `json:"day"`
}
