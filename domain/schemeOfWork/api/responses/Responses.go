package responses

type (
	SchemeOfWorkSubtopicInfo struct {
		Number             string `json:"number"`
		Code               int64  `json:"code"`
		SortCode           int64  `json:"sort_code"`
		Name               string `json:"name"`
		TopicName          string `json:"topic_name"`
		StrandName         string `json:"strand_name"`
		EstimatedTime      int64  `json:"estimated_time"`
		SchemeOfWorkNumber string `json:"scheme_of_work_number"`
		SubTopicNumber     string `json:"sub_topic_number"`
		UserNumber         string `json:"user_number"`
		ConceptCount       int64  `json:"concept_count"`
	}

	SchemeOfWorkResponseInfo struct {
		Number                   string                     `json:"number"`
		Name                     string                     `json:"name"`
		Status                   string                     `json:"status"`
		ClassNumber              string                     `json:"class_number"`
		ClassName                string                     `json:"class_name"`
		SubjectNumber            string                     `json:"subject_number"`
		SubjectName              string                     `json:"subject_name"`
		UserNumber               string                     `json:"user_number"`
		AcademicYearNumber       string                     `json:"academic_year_number"`
		AcademicYearName         string                     `json:"academic_year_name"`
		SchemeOfWorkSubtopicInfo []SchemeOfWorkSubtopicInfo `json:"subtopic"`
		SchoolName               string                     `json:"school_name"`
		SchoolBadge              string                     `json:"school_badge"`
		CampusName               string                     `json:"campus_name"`
		LevelName                string                     `json:"level_name"`
	}
)

type Report struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
