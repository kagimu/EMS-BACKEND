package requests

import "errors"

type TaskRequest struct {
	DatePerformed  string   `json:"date_performed"`
	TimeFrame      string   `json:"time_frame"`
	InputsUsed     string   `json:"inputs_used"`
	Outputs        string   `json:"outputs"`
	PeopleInvolved []string `json:"people_involved"`
	Title          string   `json:"title" binding:"required"`
	Description    string   `json:"description" binding:"required"`
	Status         string   `json:"status" binding:"required"`
	StaffNumber    string   `json:"staff_number"`
}

func (t *TaskRequest) Validate() error {
	if t.DatePerformed == "" {
		return errors.New("date performed is required")
	}
	if t.TimeFrame == "" {
		return errors.New("time frame is required")
	}
	if t.InputsUsed == "" {
		return errors.New("inputs used is required")
	}
	if len(t.PeopleInvolved) == 0 {
		return errors.New("people involved is required")
	}
	if t.Title == "" {
		return errors.New("title is required")
	}
	if t.Description == "" {
		return errors.New("description is required")
	}
	if t.Status == "" {
		return errors.New("status is required")
	}

	return nil
}
