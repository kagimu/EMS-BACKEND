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
