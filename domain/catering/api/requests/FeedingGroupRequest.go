package requests

import "errors"

type FeedingGroupRequest struct {
	MinFeedingGroupRequest
	SchoolNumber string `json:"school_number"`
}

type UpdateFeedingGroupRequest struct {
	MinFeedingGroupRequest
	FeedingGroupNumber string `json:"feeding_group_number"`
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

func (f *UpdateFeedingGroupRequest) Validate() error {
	if f.Name == "" {
		return errors.New("name is required")
	}
	if f.Status == "" {
		return errors.New("status is required")
	}
	if f.FeedingGroupNumber == "" {
		return errors.New("set the group number you are updating")
	}

	return nil
}

type MinFeedingGroupRequest struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type FeedingGroupCalendarRequest struct {
	FeedingGroups  []string `json:"feeding_groups"`
	CalendarEvents []string `json:"calendar_events"`
}

func (f *FeedingGroupCalendarRequest) Validate() error {

	if f.FeedingGroups == nil || len(f.FeedingGroups) == 0 {
		return errors.New("select at least one group")
	}

	if f.CalendarEvents == nil || len(f.CalendarEvents) == 0 {
		return errors.New("select at least one calendar")
	}

	return nil
}
