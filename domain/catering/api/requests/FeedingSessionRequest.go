package requests

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
)

type FeedingSessionRequest struct {
	Name         string `json:"name"`
	Code         uint64 `json:"code"`
	Status       string `json:"status"`
	SchoolNumber string `json:"school_number"`
}

func (r *FeedingSessionRequest) Validate() error {
	if r.Name == "" {
		return errors.New("provide name of session")
	}
	if r.Code <= 0 {
		return errors.New("provide the sorting code")
	}
	if r.Status == "" {
		return errors.New("choose the status of session")
	}
	if r.Status != constants.StatusActive && r.Status != constants.StatusInactive {
		return errors.New("the status of this session must be one of [ACTIVE, INACTIVE]")
	}
	if r.SchoolNumber == "" {
		return errors.New("choose the school and set the number")
	}
	return nil
}
