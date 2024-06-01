package requests

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
)

type UpdateFeedingSessionRequest struct {
	FeedingSessionRequest
	FeedingSessionNumber string `json:"feeding_session_number"`
}

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

func (r *UpdateFeedingSessionRequest) Validate() error {

	err := r.FeedingSessionRequest.Validate()
	if err != nil {
		return err
	}

	if r.FeedingSessionNumber == "" {
		return errors.New("choose feeding session to update")
	}

	return nil
}
