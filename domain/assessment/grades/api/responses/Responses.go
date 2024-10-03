package responses

type NameOfAggregates struct {
	Number       string `json:"number"`
	Name         string `json:"name"`
	ShortName    string `json:"short_name"`
	SchoolNumber string `json:"school_number"`
	Code         int64  `json:"code"`
	GradeGroup   string `json:"grade_group"`
	Status       string `json:"status"`
	GroupNumber  string `json:"group_number"`
}
