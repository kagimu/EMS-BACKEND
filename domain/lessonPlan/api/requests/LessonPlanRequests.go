package requests

import (
	"errors"
)

type Steps struct {
	Number string `json:"number"`
	Code   int64  `json:"code"`
	Step   string `json:"step"`
}

type LessonPlanConceptInfo struct {
	Number                  string  `json:"number"`
	Code                    int64   `json:"code"`
	LessonPlanConceptNumber string  `json:"lesson_plan_concept_number"`
	StepList                []Steps `json:"steps"`
}

type LessonPlanInfo struct {
	SchemeSubTopicNumber  string                  `json:"scheme_sub_topic_number"`
	LessonPlanNumber      string                  `json:"lesson_plan_number"`
	LessonPlanConceptInfo []LessonPlanConceptInfo `json:"concepts"`
}

func (p *LessonPlanInfo) Validate() error {

	if p.SchemeSubTopicNumber == "" {
		return errors.New("scheme is required")
	}

	if len(p.LessonPlanConceptInfo) == 0 {
		return errors.New("concepts are required")
	}

	return nil
}
