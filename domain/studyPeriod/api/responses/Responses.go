package responses

type (
	StudyPeriodList struct {
		Number               string `json:"number"`
		Code                 string `json:"code"`
		StudyPeriodName      string `json:"study_period_name"`
		StudyPeriodStartDate string `json:"study_period_start_date"`
		StudyPeriodEndDate   string `json:"study_period_end_date"`
		AcademicYearNumber   string `json:"academic_year_number"`
		AcademicYearName     string `json:"academic_year_name"`
		UserNumber           string `json:"user_number"`
		LevelNumber          string `json:"level_number"`
	}
)
