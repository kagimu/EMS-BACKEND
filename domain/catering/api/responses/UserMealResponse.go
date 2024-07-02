package responses

type UserMealResponse struct {
	UserNumber           string               `json:"user_number"`
	Name                 string               `json:"name"`
	FeedingGroupName     string               `json:"feeding_group_name"`
	FeedingGroupNumber   string               `json:"feeding_group_number"`
	Meals                []SingleUserMealItem `json:"meals"`
	FeedingSessionNumber string               `json:"feeding_session_number"`
	Served               bool                 `json:"served"`
}
