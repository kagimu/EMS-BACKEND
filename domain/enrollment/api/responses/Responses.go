package responses

type EnrolledStudent struct {
	Number            string  `json:"number"`
	StudentSurname    string  `json:"student_surname"`
	StudentSecondName string  `json:"student_second_name"`
	StudentOtherName  string  `json:"student_other_name"`
	Score             float64 `json:"score"`
	Comment           string  `json:"comment"`
	Class             string  `json:"class"`
}
