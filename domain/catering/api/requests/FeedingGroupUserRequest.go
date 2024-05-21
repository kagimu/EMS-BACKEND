package requests

import "errors"

type FeedingGroupUserRequest struct {
	FeedingGroupNumber string   `json:"feeding_group_number"`
	Status             string   `json:"status"`
	CampusNumber       string   `json:"campus_number"`
	UserNumber         []string `json:"user_number"`
}

func (r *FeedingGroupUserRequest) Validate() error {
	if r.FeedingGroupNumber == "" {
		return errors.New("feeding group number is required")
	}
	if r.Status == "" {
		return errors.New("status is required")
	}
	if len(r.UserNumber) == 0 {
		return errors.New("at least one user is required")
	}
	if r.CampusNumber == "" {
		return errors.New("choose campus")
	}

	return nil
}
