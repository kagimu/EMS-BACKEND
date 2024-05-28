package requests

import "errors"

type CourseRequest struct {
	MinCourseRequest
}

type UpdateCourseRequest struct {
	MinCourseRequest
	MealCourseNumber string `json:"meal_course_number"`
}

func (r *CourseRequest) Validate() error {
	if r.Name == "" {
		return errors.New("provide meal course name")
	}

	if r.Code <= 0 {
		return errors.New("provide meal course code")
	}

	return nil
}

func (r *UpdateCourseRequest) Validate() error {

	if r.Name == "" {
		return errors.New("provide meal course name to update")
	}

	if r.Code <= 0 {
		return errors.New("provide meal course code to update")
	}

	return nil

}

type MinCourseRequest struct {
	Name string `json:"name"`
	Code int    `json:"code"`
}
