package requests

import "errors"

type ActivateMealCourseRequest struct {
	MealCourseNumber     string   `json:"meal_course_number"`
	FeedingSessionNumber string   `json:"feeding_session_number"`
	Campuses             []string `json:"campuses"`
	Days                 []string `json:"days"`
}

func (req *ActivateMealCourseRequest) Validate() error {

	if req.MealCourseNumber == "" {
		return errors.New("meal course is required")
	}

	if req.FeedingSessionNumber == "" {
		return errors.New("feeding session is required")
	}

	if req.Campuses == nil || len(req.Campuses) == 0 {
		return errors.New("choose at least one campus")
	}

	if req.Days == nil || len(req.Days) == 0 {
		return errors.New("choose at least one day")
	}

	return nil
}
