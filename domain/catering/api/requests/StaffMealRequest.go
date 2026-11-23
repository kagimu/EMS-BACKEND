package requests

import "errors"

type StaffMealRequest struct {
	StaffNumber          string  `json:"staff_number"`
	FeedingSessionNumber string  `json:"feeding_session_number"`
	MenuItemNumber       string  `json:"menu_item_number"`
	Date                 string  `json:"date"`
	Time                 string  `json:"time"`
	Quantity             float64 `json:"quantity"`
	MealCourseNumber     string  `json:"meal_course_number"`
	CampusNumber         string  `json:"campus_number"`
}

func (m *StaffMealRequest) Validate() error {
	if m.FeedingSessionNumber == "" {
		return errors.New("feeding session number is required")
	}
	if m.StaffNumber == "" {
		return errors.New("staff number is required")
	}
	if m.Date == "" {
		return errors.New("date is required")
	}
	if m.Time == "" {
		return errors.New("time is required")
	}
	if m.MealCourseNumber == "" {
		return errors.New("choose meal being served")
	}

	return nil
}
