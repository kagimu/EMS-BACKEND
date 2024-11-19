package requests

import (
	"errors"
)

type CurriculumDeliveryRequest struct {
	StepNumbers   []LessonStep                       `json:"step_numbers"`
	StreamNumber  string                             `json:"stream_number"`
	Date          string                             `json:"date"`
	Time          string                             `json:"time"`
	TeacherNumber string                             `json:"teacher_number"`
	Comment       string                             `json:"comment"`
	Students      []StudentCurriculumDeliveryRequest `json:"students"`
}
type LessonStep struct {
	StepNumber string `json:"step_number"`
	Comment    string `json:"comment"`
	Time       string `json:"time"`
}

type UpdateCurriculumDeliveryRequest struct {
	Number string `json:"number"`
	CurriculumDeliveryRequest
}

type StudentCurriculumDeliveryRequest struct {
	StudentNumber string `json:"student_number"`
	Comment       string `json:"comment"`
}

func (cd *CurriculumDeliveryRequest) Validate() error {
	if len(cd.Students) == 0 {
		return errors.New("please select at-least one students that attended the session")
	}

	if cd.Date == "" {
		return errors.New("please select a date")
	}

	if cd.Time == "" {
		return errors.New("please select a time")
	}

	if cd.TeacherNumber == "" {
		return errors.New("please select a teacher")
	}

	if cd.StreamNumber == "" {
		return errors.New("please select the stream the lesson was delivered")
	}

	if len(cd.StepNumbers) == 0 {
		return errors.New("please select the steps that were delivered")
	}

	return nil

}
