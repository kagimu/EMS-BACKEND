package responses

type (
	LessonPlanResponseInfo struct {
		Number                  string                `json:"number"`
		LessonPlanConceptNumber string                `json:"lesson_plan_concept_number"`
		ConceptName             string                `json:"concept_name"`
		SchemeSubTopicNumber    string                `json:"scheme_sub_topic_number"`
		Code                    int64                 `json:"code"`
		UserNumber              string                `json:"user_number"`
		LessonPlanSteps         []LessonPlanStepsInfo `json:"steps"`
		//SchemeOfWorkSubtopicInfo []SchemeOfWorkSubtopicInfo `json:"subtopic"`
	}
	LessonPlanStepsInfo struct {
		Number           string `json:"number"`
		Code             int64  `json:"code"`
		Step             string `json:"step"`
		LessonPlanNumber string `json:"lesson_plan_number"`
	}

	LessonPlanReferenceResponse struct {
		LessonPlanNumber       string `json:"lesson_plan_number"`
		LearningMaterialNumber string `json:"learning_material_number"`
		Name                   string `json:"name"`
		Code                   int64  `json:"code"`
	}
)
