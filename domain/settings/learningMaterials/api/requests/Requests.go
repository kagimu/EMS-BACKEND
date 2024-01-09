package requests

import "errors"

type MaterialSubTypeRequest struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Category string `json:"category" form:"category" binding:"required"`
	Code     string `json:"code" form:"code" binding:"required"`
}

func (mt *MaterialSubTypeRequest) Validate() error {

	if mt.Name == "" {
		return errors.New("name is required")
	}

	if mt.Category == "" {
		return errors.New("category is required")
	}

	if mt.Code == "" {
		return errors.New("code is required")
	}

	return nil
}
