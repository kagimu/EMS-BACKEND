package api

import "errors"

type SourceOfSchoolRequest struct {
	Name string `json:"name"`
}

func (r SourceOfSchoolRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
