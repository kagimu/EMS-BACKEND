package requests

import "errors"

type FeedingSessionRequest struct {
	Name         string `json:"name"`
	Code         uint64 `json:"code"`
	Status       string `json:"status"`
	CampusNumber string `json:"campus_number"`
}

func (r *FeedingSessionRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Code == 0 {
		return errors.New("code is required")
	}
	if r.Status == "" {
		return errors.New("status is required")
	}
	return nil
}
