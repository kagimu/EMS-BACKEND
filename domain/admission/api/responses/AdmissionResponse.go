package responses

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/api/responses"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/api/responses"
	responses3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/api/responses"
	responses5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/api/responses"
	responses4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
)

type AdmissionResponse struct {
	AdmissionNumber string                     `json:"admission_number"`
	Status          string                     `json:"status"`
	Student         responses.StudentResponse  `json:"student"`
	CampusResponse  responses2.CampusResponse  `json:"campus"`
	Stream          responses3.StreamResponse  `json:"stream"`
	StudyPeriod     responses5.StudyPeriodList `json:"studyPeriod"`
	Class           responses4.ClassResponse   `json:"class"`
}
