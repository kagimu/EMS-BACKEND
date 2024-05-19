package requests

import "errors"

type MealCourseToFeedingGroupRequest struct {
	FeedingGroupNumber string   `json:"feeding_group_number"`
	Days               []string `json:"days"`
	MealCourses        []string `json:"meal_courses"`
}

func (f *MealCourseToFeedingGroupRequest) Validate() error {

	if f.FeedingGroupNumber == "" {
		return errors.New("choose the feeding group")
	}

	if f.Days == nil || len(f.Days) == 0 {
		return errors.New("choose at least one day")
	}

	if f.MealCourses == nil || len(f.MealCourses) == 0 {
		return errors.New("choose at least one meal course")
	}

	return nil
}
