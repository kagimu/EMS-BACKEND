package requests

import (
	"errors"
	"time"
)

type FeedingGroupPerDayRequest struct {
	FeedingGroupNumber string       `json:"feeding_group_number"`
	Day                time.Weekday `json:"day"`
}

func (f *FeedingGroupPerDayRequest) Validate() error {
	if f.FeedingGroupNumber == "" {
		return errors.New("feeding group number is required")
	}
	if f.Day == 0 {
		return errors.New("day is required")
	}
	return nil
}
