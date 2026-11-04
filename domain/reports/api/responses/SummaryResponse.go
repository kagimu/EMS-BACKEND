package responses

type KPIValue struct {
	Value float64 `json:"value"`
	Trend float64 `json:"trend"`
}

type LevelCount struct {
	Level string `json:"level"`
	Count int    `json:"count"`
}

type GenderSplit struct {
	Male   int `json:"male"`
	Female int `json:"female"`
}

type MonthlyCount struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}

type SummaryKPIs struct {
	TotalStudents     KPIValue `json:"total_students"`
	NewThisTerm       KPIValue `json:"new_this_term"`
	AttendanceRate    KPIValue `json:"attendance_rate"`
	PendingAdmissions KPIValue `json:"pending_admissions"`
}

type SummaryResponse struct {
	KPIs              SummaryKPIs    `json:"kpis"`
	EnrollmentByLevel []LevelCount   `json:"enrollment_by_level"`
	GenderSplit       GenderSplit    `json:"gender_split"`
	MonthlyEnrollment []MonthlyCount `json:"monthly_enrollment"`
}
