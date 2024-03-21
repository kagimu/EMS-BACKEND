package requests

import "errors"

type ProgressReportRequest struct {
	ClassTeacher      string `json:"class_teacher"`
	SubjectTeacher    string `json:"subject_teacher"`
	StudentNumber     string `json:"student_number"`
	StudyPeriodNumber string `json:"study_period_number"`
	Comment           string `json:"comment"`
}

func (r *ProgressReportRequest) Validate() error {
	if r.ClassTeacher == "" {
		return errors.New("class teacher is required")
	}
	if r.SubjectTeacher == "" {
		return errors.New("subject teacher is required")
	}
	if r.StudentNumber == "" {
		return errors.New("student number is required")
	}
	if r.StudyPeriodNumber == "" {
		return errors.New("study period number is required")
	}

	return nil
}
