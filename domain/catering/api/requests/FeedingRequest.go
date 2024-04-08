package requests

import "errors"

type FeedingRequest struct {
	FeedingSessionNumber string `json:"feeding_session_number"`
	StaffNumber          string `json:"staff_number"`
	Date                 string `json:"date"`
	StockingNumber       string `json:"stocking_number"`
}

func (r *FeedingRequest) Validate() error {
	if r.FeedingSessionNumber == "" {
		return errors.New("feeding session number is required")
	}
	if r.Date == "" {
		return errors.New("date is required")
	}
	if r.StaffNumber == "" {
		return errors.New("staff number is required")
	}
	return nil
}
