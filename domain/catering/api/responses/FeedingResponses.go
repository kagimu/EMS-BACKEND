package responses

import "time"

type FeedingGroupResponse struct {
	Number       string `json:"number"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	SchoolNumber string `json:"school_number"`
	SchoolName   string `json:"school_name"`
	UserName     string `json:"user_name"`
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
	Number              string  `json:"number"`
	Name                string  `json:"name"`
	Code                int64   `json:"code"`
	Type                string  `json:"type"`
	SchoolName          string  `json:"school_name"`
	SchoolNumber        string  `json:"school_number"`
	UnitOfMeasureNumber string  `json:"unit_of_measure_number"`
	UnitOfMeasureName   string  `json:"unit_of_measure_name"`
	Parts               []Part  `json:"parts"`
	CurrentQuantity     float64 `json:"current_quantity"`
	UnitCost            float64 `json:"unit_cost"`
}

type Part struct {
	Name                      string  `json:"name"`
	Quantity                  float64 `json:"quantity"`
	Code                      int64   `json:"code"`
	Number                    string  `json:"number"`
	UnitOfMeasureAbbreviation string  `json:"unit_of_measure"`
}

type StockingResponse struct {
	Number         string  `json:"name"`
	Quantity       float64 `json:"quantity"`
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

type MealResponse struct {
	Number               string      `json:"number"`
	FeedingSessionNumber string      `json:"feeding_session_number"`
	FeedingSessionName   string      `json:"feeding_session_name"`
	FeedingSessionCode   uint64      `json:"feeding_session_code"`
	FeedingSessionStatus string      `json:"feeding_session_status"`
	StockingNumber       string      `json:"stocking_number"`
	Quantity             float64     `json:"quantity"`
	Date                 time.Time   `json:"date"`
	Time                 time.Time   `json:"time"`
	MenuItemNumber       string      `json:"menu_item_number"`
	MenuItemName         string      `json:"menu_item_name"`
	StaffIDNumber        string      `json:"staff_id_number"`
	UserMeals            []UserMeals `json:"user_meals"`
	JobTitle             string      `json:"job_title"`
}

type UserMeals struct {
	UserName     string  `json:"user_name"`
	MenuItemName string  `json:"menu_item_name"`
	Quantity     float64 `json:"quantity"`
	MenuItemType string  `json:"menu_item_type"`
	MenuItemCode uint64  `json:"menu_item_code"`
}
