package responses

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
