package requests

import (
	"errors"
)

type CurriculumDeliveryRequest struct {
	Session       []LessonStep `json:"session"`
	StreamNumber  string       `json:"stream_number"`
	Date          string       `json:"date"`
	Time          string       `json:"time"`
	TeacherNumber string       `json:"teacher_number"`
	Comment       string       `json:"comment"`
}
type LessonStep struct {
	StepNumber string                             `json:"step_number"`
	Comment    string                             `json:"comment"`
	Time       string                             `json:"time"`
	Students   []StudentCurriculumDeliveryRequest `json:"students"`
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

	for _, step := range cd.Session {
		if step.StepNumber == "" {
			return errors.New("please select a step")
		}

		if step.Time == "" {
			return errors.New("please select a time")
		}

		if len(step.Students) == 0 {
			return errors.New("please select students for the step")
		}
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

	if len(cd.Session) == 0 {
		return errors.New("please select the steps that were delivered")
	}

	return nil

}
