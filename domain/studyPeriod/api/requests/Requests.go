package requests

type StudyPeriodDelRequest struct {
	Number string `json:"number"`
}

type StudyPeriodInfo struct {
	Code                 string `json:"code"`
	StudyPeriodName      string `json:"study_period_name"`
	AcademicYearNumber   string `json:"academic_year_number"`
	StudyPeriodStartDate string `json:"study_period_start_date"`
	StudyPeriodEndDate   string `json:"study_period_end_date"`
	StudyPeriodNumber    string `json:"study_period_number"`
}
