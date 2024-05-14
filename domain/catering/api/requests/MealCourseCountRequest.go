package requests

import (
	"errors"
	"strconv"
)

type MealCourseCountRequest struct {
	SchoolNumber string          `json:"school_number"`
	Count        int             `json:"count"`
	Courses      []CourseRequest `json:"courses"`
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

	return nil
}
