package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"

type (
	StudyPeriodList struct {
		Number               string `json:"number"`
		Code                 int64  `json:"code"`
		StudyPeriodName      string `json:"study_period_name"`
		StudyPeriodStartDate string `json:"study_period_start_date"`
		StudyPeriodEndDate   string `json:"study_period_end_date"`
		UserNumber           string `json:"user_number"`
		LevelNumber          string `json:"level_number"`
		AcademicYearNumber   string `json:"academic_year_number"`
		AcademicYearName     string `json:"academic_year_name"`
	}

	StudyPeriodNameResponse struct {
		DateAdded             string `json:"date_added"`
		Name                  string `json:"name"`
		StudyPeriodNameNumber string `json:"study_period_name_number"`
	}

	StudyPeriodListDetail struct {
		StudyPeriodList
		Classes []responses.ClassResponse `json:"classes"`
	}
)
