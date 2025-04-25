package responses

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/api/responses"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/api/responses"
	responses3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/api/responses"
	responses5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/api/responses"
	responses4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
)

type AdmissionResponse struct {
	AdmissionNumber    string                     `json:"admission_number"`
	Status             string                     `json:"status"`
	Student            responses.StudentResponse  `json:"student"`
	CampusResponse     responses2.CampusResponse  `json:"campus"`
	Stream             responses3.StreamResponse  `json:"stream"`
	StudyPeriod        responses5.StudyPeriodList `json:"studyPeriod"`
	Class              responses4.ClassResponse   `json:"class"`
	SchoolNumber       string                     `json:"school_number"`
	CampusNumber       string                     `json:"campus_number"`
	CurriculumNumber   string                     `json:"curriculum_number"`
	LevelNumber        string                     `json:"level_number"`
	AcademicYearNumber string                     `json:"academic_year_number"`
	StudyPeriodNumber  string                     `json:"study_period_number"`
	CampusClassNumber  string                     `json:"campus_class_number"`
	StreamNumber       string                     `json:"stream_number"`
}
