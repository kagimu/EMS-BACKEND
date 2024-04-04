package responses

type EnrolledStudent struct {
	Number            string `json:"number"`
	StudentSurname    string `json:"student_surname"`
	StudentSecondName string `json:"student_second_name"`
	StudentOtherName  string `json:"student_other_name"`
}
