package requests

import (
	"errors"
	"fmt"
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

	for index, lessonInfo := range p.LessonPlanConceptInfo {
		position := index + 1
		if lessonInfo.LessonPlanConceptNumber == "" {
			return errors.New(fmt.Sprintf("lesson_plan_concept_number is required, at position %d", position))
		}
	}

	hasDuplicateCode, code := hasDuplicateCodes(p.LessonPlanConceptInfo)
	if hasDuplicateCode {
		return errors.New(fmt.Sprintf("Duplicate code entry '%d' in concepts", code))
	}

	hasDuplicateConcept, _ := hasDuplicateConcepts(p.LessonPlanConceptInfo)
	if hasDuplicateConcept {
		return errors.New(fmt.Sprintf("Duplicate concept entry in the selected concepts"))
	}

	return nil
}

func hasDuplicateCodes(subtopics []LessonPlanConceptInfo) (bool, int64) {
	seen := make(map[int64]bool)
	for _, s := range subtopics {
		if seen[s.Code] {
			return true, s.Code
		}
		seen[s.Code] = true
	}
	return false, 0
}

func hasDuplicateConcepts(subtopics []LessonPlanConceptInfo) (bool, string) {
	seen := make(map[string]bool)
	for _, s := range subtopics {
		if seen[s.LessonPlanConceptNumber] {
			return true, s.LessonPlanConceptNumber
		}
		seen[s.LessonPlanConceptNumber] = true
	}
	return false, ""
}
