package responses

import (
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
}
