package requests

import "errors"

type FeedingGroupStaffRequest struct {
	FeedingGroupNumber string   `json:"feeding_group_number"`
	Status             string   `json:"status"`
	CampusNumber       string   `json:"campus_number"`
	StaffNumber        []string `json:"staff_number"`
}

func (r *FeedingGroupStaffRequest) Validate() error {
	if r.FeedingGroupNumber == "" {
		return errors.New("feeding group number is required")
	}
	if r.Status == "" {
		return errors.New("status is required")
	}
	if len(r.StaffNumber) == 0 {
		return errors.New("at least one staff member is required")
	}
	return nil
}
