package requests

import (
	"errors"
	"strconv"
)

type MealCourseCountRequest struct {
	MinMealCourseCountRequest
	Courses []CourseRequest `json:"courses"`
}

type UpdateMealCourseCountRequest struct {
	MinMealCourseCountRequest
	Courses []UpdateCourseRequest `json:"courses"`
}

func (mealCourseCountRequest *MealCourseCountRequest) Validate() error {

	if mealCourseCountRequest.SchoolNumber == "" {
		return errors.New("provide school number")
	}

	if mealCourseCountRequest.Count < 1 {
		return errors.New("provide meal course count")
	}

	if len(mealCourseCountRequest.Courses) != mealCourseCountRequest.Count {
		return errors.New("meal courses defined should be " + strconv.Itoa(mealCourseCountRequest.Count))
	}

	for _, course := range mealCourseCountRequest.Courses {
		if err := course.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (mealCourseCountRequest *UpdateMealCourseCountRequest) Validate() error {

	if mealCourseCountRequest.SchoolNumber == "" {
		return errors.New("provide school number")
	}

	if mealCourseCountRequest.Count < 1 {
		return errors.New("provide meal course count")
	}

	if len(mealCourseCountRequest.Courses) != mealCourseCountRequest.Count {
		return errors.New("meal courses defined should be " + strconv.Itoa(mealCourseCountRequest.Count))
	}

	for _, course := range mealCourseCountRequest.Courses {
		if err := course.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type MinMealCourseCountRequest struct {
	SchoolNumber string `json:"school_number"`
	Count        int    `json:"count"`
}
