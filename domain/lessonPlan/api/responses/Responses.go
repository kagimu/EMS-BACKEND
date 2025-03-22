package responses

type (
	LessonPlanData struct {
		Number               string                        `json:"number"`
		SchemeSubTopicNumber string                        `json:"scheme_sub_topic_number"`
		Status               string                        `json:"status"`
		UserNumber           string                        `json:"user_number"`
		LessonPlanConcepts   []LessonPlanConceptData       `json:"concepts"`
		References           []LessonPlanReferenceResponse `json:"references"`
	}

	LessonPlanConceptData struct {
		Number           string                `json:"number"`
		Code             int64                 `json:"code"`
		ConceptNumber    string                `json:"lesson_plan_concept_number"`
		ConceptName      string                `json:"concept_name"`
		LessonPlanNumber string                `json:"lesson_plan_number"`
		UserNumber       string                `json:"user_number"`
		LessonPlanSteps  []LessonPlanStepsData `json:"steps"`
	}

	LessonPlanStepsData struct {
		Number                  string `json:"number"`
		Code                    int64  `json:"code"`
		Step                    string `json:"step"`
		LessonPlanConceptNumber string `json:"lesson_plan_concept_number"`
	}

	LessonPlanReferenceResponse struct {
		LessonPlanNumber          string `json:"lesson_plan_number"`
		LearningMaterialNumber    string `json:"learning_material_number"`
		Name                      string `json:"name"`
		Code                      int64  `json:"code"`
		LessonPlanReferenceNumber string `json:"lesson_plan_reference_number"`
	}
)
