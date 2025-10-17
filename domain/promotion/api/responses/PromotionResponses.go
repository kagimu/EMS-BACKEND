package responses

import (
	responses7 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/responses"
	responses6 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/api/responses"
	responses4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/api/responses"
	responses5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/api/responses"
	responses3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/api/responses"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
)

type PromotionResponse struct {
	PromotionSettingsNumber  string                         `json:"promotion_settings_number"`
	Name                     string                         `json:"name"`
	PreviousClass            responses.ClassResponse        `json:"previous_class"`
	NextClass                responses.ClassResponse        `json:"next_class"`
	StudyPeriod              responses2.StudyPeriodResponse `json:"study_period"`
	CurrentStudyPeriodNumber string                         `json:"current_study_period_number"`
	Status                   string                         `json:"status"`
	Type                     string                         `json:"type"`
	Curriculum               responses3.CurriculumResponse  `json:"curriculum"`
	Campus                   responses4.CampusResponse      `json:"campus"`
	SchoolName               string                         `json:"school_name"`
	SchoolNumber             string                         `json:"school_number"`
	CampusCurriculumNumber   string                         `json:"campus_curriculum_number"`
	CampusLevelNumber        string                         `json:"campus_level_number"`
}

type PromotionDetailedResponse struct {
	PromotionResponse
	NextCampusClassNumber     string `json:"next_campus_class_number"`
	PreviousCampusClassNumber string `json:"previous_campus_class_number"`
}

type PromotionSettingGradeResponse struct {
	AssessmentAreas []responses5.AssessmentAreaResponse      `json:"assessment_areas"`
	GradeScales     []responses6.AssessmentGradeItemResponse `json:"grade_scales"`
	AggregateScales []responses7.AggregateScaleResponse      `json:"aggregate_scales"`
}

type PromotionSettingValueResponse struct {
	PromotionDetailedResponse
	AssessmentAreas       []PromotionAssessmentAreaValueResponse `json:"assessment_areas"`
	AssessmentGradeName   string                                 `json:"assessment_grade_name"`
	AssessmentGradeNumber string                                 `json:"assessment_grade_number"`
	MinimumAggregates     float64                                `json:"minimum_aggregates"`
	MinimumScore          float64                                `json:"minimum_score"`
}

type PromotionAssessmentAreaValueResponse struct {
	AssessmentArea       string `json:"assessment_area"`
	AggregateScale       string `json:"aggregate_scale"`
	Number               string `json:"number"`
	AssessmentAreaNumber string `json:"assessment_area_number"`
	AggregateScaleNumber string `json:"aggregate_scale_number"`
}
