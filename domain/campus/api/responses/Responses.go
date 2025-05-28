package responses

import (
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/api/responses"
	"time"
)

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
		Number                 string `json:"number"`
		SchoolCurriculumNumber string `json:"school_curriculum_number"`
		Name                   string `json:"name"`
		CurriculumNumber       string `json:"curriculum_number"`
	}

	CampusLevelDataList struct {
		Number                string                        `json:"number"`
		CurriculumLevelNumber string                        `json:"curriculum_level_number"`
		Code                  int64                         `json:"code"`
		Name                  string                        `json:"name"`
		Classes               []responses.CampusClassInfo   `json:"classes"`
		AcademicYear          []responses2.AcademicYearList `json:"academicYear"`
	}

	CampusLevelList struct {
		Number              string `json:"number"`
		LearningLevelNumber string `json:"learning_level_number"`
		LevelName           string `json:"level_name"`
		CampusNumber        string `json:"campus_number"`
	}

	CampusResponse struct {
		Number           string    `json:"number"`
		Name             string    `json:"name"`
		Phone            string    `json:"phone"`
		Email            string    `json:"email"`
		Location         string    `json:"location"`
		SchoolNumber     string    `json:"school_number"`
		CommencementDate time.Time `json:"commencement_date" `
	}
)
