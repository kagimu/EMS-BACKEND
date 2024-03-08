package requests

import "errors"

type JobTitleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (j *JobTitleRequest) Validate() error {
	if j.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
