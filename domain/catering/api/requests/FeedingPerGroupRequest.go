package requests

import "errors"

type FeedingPerGroupRequest struct {
	FeedingGroupNumber   string `json:"feeding_group_number"`
	Day                  string `json:"day"`
	FeedingSessionNumber string `json:"feeding_session_number"`
}

func (f *FeedingPerGroupRequest) Validate() error {
	if f.FeedingGroupNumber == "" {
		return errors.New("feeding group number is required")
	}
	if f.Day == "" {
		return errors.New("day is required")
	}
	if f.FeedingSessionNumber == "" {
		return errors.New("feeding session is required")
	}

	return nil
}
