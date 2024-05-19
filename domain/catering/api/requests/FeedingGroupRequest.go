package requests

import "errors"

type FeedingGroupRequest struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	SchoolNumber string `json:"school_number"`
}

func (f *FeedingGroupRequest) Validate() error {
	if f.Name == "" {
		return errors.New("name is required")
	}
	if f.Status == "" {
		return errors.New("status is required")
	}
	if f.SchoolNumber == "" {
		return errors.New("choose the school")
	}

	return nil
}
