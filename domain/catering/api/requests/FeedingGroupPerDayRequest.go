package requests

import (
	"errors"
)

type FeedingGroupPerDayRequest struct {
	FeedingGroupNumber string `json:"feeding_group_number"`
	Day                string `json:"day"`
}

func (f *FeedingGroupPerDayRequest) Validate() error {
	if f.FeedingGroupNumber == "" {
		return errors.New("feeding group number is required")
	}
	if f.Day == "" {
		return errors.New("day is required")
	}
	return nil
}
