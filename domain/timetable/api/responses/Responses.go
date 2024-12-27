package responses

import "time"

type Subjects struct {
	Number string `json:"number"`
	Name   string `json:"name"`
}

type ClassResponse struct {
	Number       string `json:"number"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	LevelName    string `json:"level_name"`
	LevelNumber  string `json:"level_number"`
	Code         int64  `json:"code"`
}

type StudyTimetables struct {
	Name   string `json:"name"`
	Number string `json:"number"`
	Status string `json:"status"`
}

type StudyingTimeBlockResponse struct {
	Number             string    `json:"number"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	Length             uint      `json:"length"`
	AcademicYearNumber string    `json:"academic_year_number"`
	AcademicYearName   string    `json:"academic_year_name"`
	CampusName         string    `json:"campus_name"`
	LevelName          string    `json:"level_name"`
}

type StudyingTimeTableResponse struct {
	Number                     string    `json:"number"`
	StudyingTimeBlockNumber    string    `json:"studying_time_block_number"`
	StudyingTimeBlockStartTime time.Time `json:"studying_time_block_start_time"`
	StudyingTimeBlockEndTime   time.Time `json:"studying_time_block_end_time"`
	StreamName                 string    `json:"stream_name"`
	Day                        time.Time `json:"day"`
	SubjectName                string    `json:"subject_name"`
	StartTime                  time.Time `json:"start_time"`
	EndTime                    time.Time `json:"end_time"`
}

type StudyingTimeTablePeriodResponse struct {
	Number                     string    `json:"number"`
	StudyingTimeTableNumber    string    `json:"studying_time_table_number"`
	UserName                   string    `json:"user_name"`
	Comment                    string    `json:"comment"`
	LessonPlanStepName         string    `json:"lesson_plan_step_name"`
	StudyingTimeTableStartTime time.Time `json:"studying_time_table_start_time"`
	StudyingTimeTableEndTime   time.Time `json:"studying_time_table_end_time"`
}
