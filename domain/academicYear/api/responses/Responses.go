package responses

type (
	AcademicYearList struct {
		Number                string `json:"number"`
		Code                  int64  `json:"code"`
		AcademicYearName      string `json:"academic_year_name"`
		NoOfStudyPeriod       int32  `json:"no_of_study_period"`
		AcademicYearStartDate string `json:"academic_year_start_date"`
		AcademicYearEndDate   string `json:"academic_year_end_date"`
		UserNumber            string `json:"user_number"`
		LevelNumber           string `json:"level_number"`
	}

	AcademicYearListByClass struct {
		Number                string `json:"number"`
		Code                  int64  `json:"code"`
		AcademicYearName      string `json:"academic_year_name"`
		NoOfStudyPeriod       int32  `json:"no_of_study_period"`
		AcademicYearStartDate string `json:"academic_year_start_date"`
		AcademicYearEndDate   string `json:"academic_year_end_date"`
		UserNumber            string `json:"user_number"`
		LevelNumber           string `json:"level_number"`
		LevelName             string `json:"level_name"`
	}
)
