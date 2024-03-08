package requests

import "errors"

type ReportCardRequest struct {
	ClassTeacher    string `json:"class_teacher"`
	SubjectTeacher  string `json:"subject_teacher"`
	StudentNumber   string `json:"student_number"`
	AssesmentNumber string `json:"assessment_number"`
	Comment         string `json:"comment"`
}

func (r *ReportCardRequest) Validate() error {
	if r.ClassTeacher == "" {
		return errors.New("class teacher is required")
	}
	if r.SubjectTeacher == "" {
		return errors.New("subject teacher is required")
	}

	return nil
}
