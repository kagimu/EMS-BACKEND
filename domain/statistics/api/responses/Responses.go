package responses

type BasicStatistics struct {
	TotalStudents int `json:"total_students"`
	TotalTeachers int `json:"total_teachers"`
	TotalClasses  int `json:"total_classes"`
	TotalCampuses int `json:"total_campuses"`
	TotalParents  int `json:"total_parents"`
}
