package responses

type (
	CampusDataList struct {
		Number       string `json:"number"`
		Name         string `json:"name"`
		Phone        string `json:"phone"`
		Email        string `json:"email"`
		Location     string `json:"location"`
		UserNumber   string `json:"user_number"`
		SchoolNumber string `json:"school_number"`
		//SchoolCurriculumNumber string `json:"school_curriculum_number"`
		//CurriculumName         string               `json:"curriculum_name"`
		CommencementDate string               `json:"commencement_date"`
		Curriculum       []CurriculumDataList `json:"curriculum"`
	}

	CurriculumDataList struct {
		Number string `json:"number"`
		Name   string `json:"name"`
	}

	CampusLevelList struct {
		Number              string `json:"number"`
		LearningLevelNumber string `json:"learning_level_number"`
		LevelName           string `json:"level_name"`
		CampusNumber        string `json:"campus_number"`
	}
)
