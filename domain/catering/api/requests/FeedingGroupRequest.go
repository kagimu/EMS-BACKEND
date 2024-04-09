package requests

import "errors"

type FeedingGroupRequest struct {
	Name         string   `json:"name"`
	Status       string   `json:"status"`
	CampusNumber string   `json:"campus_number"`
	UserNumber   []string `json:"user_number"`
}

func (f *FeedingGroupRequest) Validate() error {
	if f.Name == "" {
		return errors.New("name is required")
	}
	if f.Status == "" {
		return errors.New("status is required")
	}
	if f.CampusNumber == "" {
		return errors.New("campus number is required")
	}
	if len(f.UserNumber) == 0 {
		return errors.New("at least one user is required")
	}

	return nil
}
