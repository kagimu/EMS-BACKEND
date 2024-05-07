package requests

import (
	"errors"
)

type FeedingSessionPerDayRequest struct {
	Day                  []string `json:"day"`
	FeedingSessionNumber string   `json:"feeding_session_number"`
}

func (f *FeedingSessionPerDayRequest) Validate() error {
	if len(f.Day) == 0 {
		return errors.New("please select at least one day of the week")
	}
	if f.FeedingSessionNumber == " " {
		return errors.New("feeding session number is required")
	}

	return nil
}
