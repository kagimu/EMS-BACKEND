package requests

import "errors"

type FeedingGroupSessionRequest struct {
	FeedingGroupNumber   string `json:"feeding_group_number"`
	Status               string `json:"status"`
	FeedingSessionNumber string `json:"feeding_session_number"`
}

func (f *FeedingGroupSessionRequest) Validate() error {
	if f.FeedingGroupNumber == "" {
		return errors.New("feeding group number is required")
	}
	if f.FeedingSessionNumber == "" {
		return errors.New("feeding session number is required")
	}

	return nil
}
