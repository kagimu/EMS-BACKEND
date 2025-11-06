package requests

import "time"

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

type StudyTimeBlockRequest struct {
	StartTime          string        `json:"start_time"`
	EndTime            string        `json:"end_time"`
	Day                string        `json:"day"`
	CampusNumber       string        `json:"campus_number"`
	Length             time.Duration `json:"length"`
	AcademicYearNumber string        `json:"academic_year_number"`
	LevelNumber        string        `json:"level_number"`
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

type TimetableTaskRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	ClassNumber string `json:"class_number" binding:"required"`
}

type UpdateTimetableTaskRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	ClassNumber string `json:"class_number"`
}

type TimetableSlotRequest struct {
	Name         string `json:"name" binding:"required"`
	Type         string `json:"type" binding:"required"`
	SchoolNumber string `json:"school_number" binding:"required"`
	LevelNumber  string `json:"level_number" binding:"required"`
	StartTime    string `json:"start_time" binding:"required"`
	EndTime      string `json:"end_time" binding:"required"`
}

type UpdateTimetableSlotRequest struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	SchoolNumber string `json:"school_number"`
	LevelNumber  string `json:"level_number"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
}
