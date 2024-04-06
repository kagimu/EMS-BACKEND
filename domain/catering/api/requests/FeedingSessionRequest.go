package requests

import "errors"

type FeedingSesionRequest struct {
	Name   string `json:"name"`
	Code   uint64 `json:"code"`
	Status string `json:"status"`
}

func (r *FeedingSesionRequest) Validate() error {
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
