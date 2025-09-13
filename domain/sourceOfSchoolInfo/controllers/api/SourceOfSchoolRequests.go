package api

import "errors"

type SourceOfSchoolRequest struct {
	Name string `json:"name"`
	Code int64  `json:"code"`
}

func (r SourceOfSchoolRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Code == 0 {
		return errors.New("code is required")
	}

	return nil
}
