package requests

import (
	"errors"
)

type FeedingSessionPerDayRequest struct {
	Day                  string `json:"day"`
	FeedingSessionNumber string `json:"feeding_session_number"`
}

func (f *FeedingSessionPerDayRequest) Validate() error {
	if f.Day == "" {
		return errors.New("week day is required")
	}
	if f.FeedingSessionNumber == " " {
		return errors.New("feeding session number is required")
	}

	return nil
}
