package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"

type GradingStandardResponse struct {
	Code        int64  `json:"code" `
	Number      string `json:"number" `
	Name        string `json:"name"`
	Description string `json:"description" `
	Status      string `json:"status"`
}

type DetailedGradingStandardResponse struct {
	GradingStandardResponse
	StandardClasses []GradingStandardClassesResponse `json:"standard_classes"`
}

type GradingStandardClassesResponse struct {
	Class responses.ClassResponse `json:"class"`
	Scale DetailedScaleResponse   `json:"scale"`
}

type ClassGradingStandardResponse struct {
	Class           responses.ClassResponse `json:"class"`
	Scale           ScaleResponse           `json:"scale"`
	GradingStandard GradingStandardResponse `json:"grading_standard"`
	Status          string                  `json:"status"`
}
