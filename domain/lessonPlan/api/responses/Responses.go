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
	/*	SchemeOfWorkSubtopicInfo struct {
		Number             string `json:"number"`
		Code               int64  `json:"code"`
		SortCode           int64  `json:"sort_code"`
		Name               string `json:"name"`
		TopicName          string `json:"topic_name"`
		StrandName         string `json:"strand_name"`
		EstimatedTime      int64  `json:"estimated_time"`
		SchemeOfWorkNumber string `json:"scheme_of_work_number"`
		SubtopicNumber     string `json:"sub_topic_number"`
		UserNumber         string `json:"user_number"`
	}*/
)
