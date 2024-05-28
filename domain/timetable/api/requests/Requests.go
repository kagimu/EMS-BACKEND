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
