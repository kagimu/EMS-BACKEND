package responses

import (
	responses4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/api/responses"
	responses3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/api/responses"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
)

type PromotionResponse struct {
	PromotionSettingsNumber string                         `json:"promotion_settings_number"`
	Name                    string                         `json:"name"`
	PreviousClass           responses.ClassResponse        `json:"previous_class"`
	NextClass               responses.ClassResponse        `json:"next_class"`
	StudyPeriod             responses2.StudyPeriodResponse `json:"study_period"`
	Status                  string                         `json:"status"`
	Type                    string                         `json:"type"`
	Curriculum              responses3.CurriculumResponse  `json:"curriculum"`
	Campus                  responses4.CampusResponse      `json:"campus"`
}

type PromotionDetailedResponse struct {
	PromotionResponse
	NextCampusClassNumber     string `json:"next_campus_class_number"`
	PreviousCampusClassNumber string `json:"previous_campus_class_number"`
}
