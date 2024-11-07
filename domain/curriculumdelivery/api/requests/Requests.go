package requests

import "errors"

type CurriculumDeliveryRequest struct {
	LessonPlanStepNumbers []string                           `json:"lesson_plan_numbers" form:"lesson_plan_step_numbers" binding:"required"`
	StreamNumber          string                             `json:"stream_number" form:"stream_number" binding:"required"`
	Date                  string                             `json:"date" form:"date" binding:"required"`
	Time                  string                             `json:"time" form:"time" binding:"required"`
	TeacherNumber         string                             `json:"teacher_number" form:"teacher_number" binding:"required"`
	Comment               string                             `json:"comment" form:"comment"`
	Students              []StudentCurriculumDeliveryRequest `json:"students"`
}

type StudentCurriculumDeliveryRequest struct {
	StudentNumber string `json:"student_number" form:"student_number" binding:"required"`
	Comment       string `json:"comment" form:"comment"`
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

	if len(cd.LessonPlanStepNumbers) < 1 {
		return errors.New("please select the lesson plan that was delivered")
	}

	return nil

}
