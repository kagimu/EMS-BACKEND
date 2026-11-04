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

type DeleteMealCourseFeedingSessionRequest struct {
	MealCourseNumber     string `json:"meal_course_number"`
	FeedingSessionNumber string `json:"feeding_session_number"`
}

func (req *DeleteMealCourseFeedingSessionRequest) Validate() error {
	if req.MealCourseNumber == "" {
		return errors.New("provide meal course number")
	}
	if req.FeedingSessionNumber == "" {
		return errors.New("provide feeding session number")
	}
	return nil
}
