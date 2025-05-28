package responses

type FeedingGroupResponse struct {
	Number       string `json:"number"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	SchoolNumber string `json:"school_number"`
	SchoolName   string `json:"school_name"`
	UserName     string `json:"user_name"`
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
	Name               string `json:"name"`
	FeedingGroupNumber string `json:"feeding_group_number"`
	FeedingGroupName   string `json:"feeding_group_name"`
	CampusName         string `json:"campus_name"`
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
	Number                      string  `json:"number"`
	Name                        string  `json:"name"`
	Code                        int64   `json:"code"`
	Type                        string  `json:"type"`
	SchoolName                  string  `json:"school_name"`
	SchoolNumber                string  `json:"school_number"`
	UnitOfMeasureNumber         string  `json:"unit_of_measure_number"`
	UnitOfMeasureName           string  `json:"unit_of_measure_name"`
	PurchaseUnitOfMeasureNumber string  `json:"purchase_unit_of_measure_number"`
	PurchaseUnitOfMeasureName   string  `json:"purchase_unit_of_measure_name"`
	Parts                       []Part  `json:"parts"`
	CurrentQuantity             float64 `json:"current_quantity"`
	PurchaseCurrentQuantity     float64 `json:"purchase_current_quantity"`
	UnitCost                    float64 `json:"unit_cost"`
	PurchaseUnitCost            float64 `json:"purchase_unit_cost"`
	ConversionToPurchaseUnit    float64 `json:"conversion_to_purchase_unit"`
}

type Part struct {
	Name                      string  `json:"name"`
	Quantity                  float64 `json:"quantity"`
	Code                      int64   `json:"code"`
	Number                    string  `json:"number"`
	UnitOfMeasureAbbreviation string  `json:"unit_of_measure"`
	UnitCost                  float64 `json:"unit_cost"`
}

type StockingResponse struct {
	Number           string  `json:"name"`
	Quantity         float64 `json:"quantity"`
	PurchaseQuantity float64 `json:"purchase_quantity"`
	PurchaseUnitName string  `json:"purchase_unit_name"`
	UserUnitName     string  `json:"user_unit_name"`
	UnitCost         float64 `json:"unit_cost"`
	PurchaseUnitCost float64 `json:"purchase_unit_cost"`
	MenuItemNumber   string  `json:"menu_item_number"`
	MenuItemName     string  `json:"menu_item_name"`
	CampusNumber     string  `json:"campus_number"`
	CampusName       string  `json:"campus_name"`
	MenuItemType     string  `json:"menu_item_type"`
	DateTime         string  `json:"date_time"`
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

type FeedingGroupPerDayResponse struct {
	Number             string `json:"number"`
	FeedingGroupNumber string `json:"feeding_group_number"`
	FeedingGroupName   string `json:"feeding_group_name"`
	Day                string `json:"day"`
}

type UserMeals struct {
	UserName     string  `json:"user_name"`
	MenuItemName string  `json:"menu_item_name"`
	Quantity     float64 `json:"quantity"`
	MenuItemType string  `json:"menu_item_type"`
	MenuItemCode uint64  `json:"menu_item_code"`
}
