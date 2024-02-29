package requests

import "errors"

type ExaminationTypeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (e *ExaminationTypeRequest) Validate() error {
	if e.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
