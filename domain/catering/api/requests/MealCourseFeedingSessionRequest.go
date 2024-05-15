package requests

import (
	"errors"
)

type MealCourseFeedingSessionRequest struct {
	MealCourseNumber string   `json:"meal_course_number"`
	FeedingSessions  []string `json:"feeding_sessions"`
}

func (mealCourseFeedingSessionRequest *MealCourseFeedingSessionRequest) Validate() error {

	if mealCourseFeedingSessionRequest.MealCourseNumber == "" {
		return errors.New("provide meal course number")
	}

	if len(mealCourseFeedingSessionRequest.FeedingSessions) < 1 {
		return errors.New("provide meal course feeding sessions")
	}

	return nil
}
