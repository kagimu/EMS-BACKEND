package responses

type BillableResponse struct {
	BillableNumber         string  `json:"billable_number"`
	Status                 string  `json:"status"`
	ClassName              string  `json:"class_name"`
	ClassNumber            string  `json:"class_number"`
	StudyPeriodNumber      string  `json:"study_period_number"`
	StudyPeriodName        string  `json:"study_period_name"`
	CampusNumber           string  `json:"campus_number"`
	CampusName             string  `json:"campus_name"`
	BillType               string  `json:"bill_type"`
	TotalAmount            float64 `json:"total_amount"`
	CampusCurriculumNumber string  `json:"campus_curriculum_number"`
	AcademicYearNumber     string  `json:"academic_year_number"`
	AcademicYearName       string  `json:"academic_year_name"`
	CampusLevelNumber      string  `json:"campus_level_number"`
}

type DetailedBillableResponse struct {
	BillableResponse
	BillableItems []InvoiceItemResponse `json:"billable_items"`
}
