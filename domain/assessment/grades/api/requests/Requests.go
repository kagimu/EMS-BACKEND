package requests

import "errors"

type GradeDataRequest struct {
	Name         string `json:"name"`
	Code         int64  `json:"code"`
	ShortName    string `json:"short_name"`
	SchoolNumber string `json:"school_number"`
}

type CreateGradeGroupRequest struct {
	Code         int64  `json:"code"`
	Name         string `json:"name"`
	SchoolNumber string `json:"school_number"`
}

type UpdateGradeGroupRequest struct {
	Number string `json:"number"`
	CreateGradeGroupRequest
}

type EditGradeDataRequest struct {
	GradeNumber string `json:"number"`
	GradeDataRequest
}

func (g *GradeDataRequest) Validate() error {
	if g.Name == "" {
		return errors.New("name is required")
	}
	if g.Code == 0 {
		return errors.New("code is required")
	}
	if g.ShortName == "" {
		return errors.New("short name is required")
	}

	return nil
}

func (g *EditGradeDataRequest) Validate() error {

	if err := g.GradeDataRequest.Validate(); err != nil {
		return err
	}

	if g.GradeNumber == "" {
		return errors.New("number is required")
	}
	return g.GradeDataRequest.Validate()
}

func (c *CreateGradeGroupRequest) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	if c.SchoolNumber == "" {
		return errors.New("school number is required")
	}

	return nil
}

func (u *UpdateGradeGroupRequest) Validate() error {
	if u.Number == "" {
		return errors.New("number is required")
	}
	if err := u.CreateGradeGroupRequest.Validate(); err != nil {
		return err
	}
	return nil
}
