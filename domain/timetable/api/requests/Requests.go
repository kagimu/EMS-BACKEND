package requests

type TimetableRequest struct {
	SubjectNumber string `json:"subject_number"`
	StreamNumber  string `json:"stream_number"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	Day           string `json:"day"`
	CampusNumber  string `json:"campus_number"`
}

type StudyTimetableRequest struct {
	StudyPeriodNumber string `json:"study_period_number"`
	Name              string `json:"name"`
}

type StudyingTimeTableRequest struct {
	StudyingTimeBlockNumber string `json:"studying_time_block_number"`
	StreamNumber            string `json:"stream_number"`
	SubjectNumber           string `json:"subject_number"`
	StartTime               string `json:"start_time"`
	EndTime                 string `json:"end_time"`
	Day                     string `json:"day"`
}

type StudyPeriodRequest struct {
	StudyingTimeTableNumber string `json:"studying_time_table_number"`
	Comment                 string `json:"comment"`
	LessonPlanStepNumber    string `json:"lesson_plan_step_number"`
}

type TimetableTaskItem struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
}

type TimetableTaskRequest struct {
	Tasks        []TimetableTaskItem `json:"tasks" binding:"required,min=1,dive"`
	ClassNumbers []string            `json:"class_numbers" binding:"required,min=1,dive,required"`
}

type UpdateTimetableTaskRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	ClassNumber string `json:"class_number"`
}

type TimetableSlotItem struct {
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type TimetableSlotRequest struct {
	Name         string              `json:"name" binding:"required"`
	Type         string              `json:"type" binding:"required"`
	SchoolNumber string              `json:"school_number" binding:"required"`
	LevelNumber  string              `json:"level_number" binding:"required"`
	Slots        []TimetableSlotItem `json:"slots" binding:"required,min=1,dive"`
}

type UpdateTimetableSlotRequest struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	SchoolNumber string `json:"school_number"`
	LevelNumber  string `json:"level_number"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
}

type StudyTimetablePeriodItem struct {
	Day        string `json:"day" binding:"required"`
	TaskNumber string `json:"task_number" binding:"required"`
	SlotNumber string `json:"slot_number" binding:"required"`
}

type CreateStudyTimetablePeriodRequest struct {
	StudyPeriodNumber string                     `json:"study_period_number"`
	StreamNumber      string                     `json:"stream_number" binding:"required"`
	Periods           []StudyTimetablePeriodItem `json:"periods"`
}

type UpdateStudyTimetablePeriodRequest struct {
	StudyTimetableNumber string `json:"study_timetable_number" binding:"required"`
	CreateStudyTimetablePeriodRequest
}
