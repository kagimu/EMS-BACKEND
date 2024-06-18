package requests

import "errors"

type GradeDataRequest struct {
	Name        string  `json:"name"`
	Code        int64   `json:"code"`
	ShortName   string  `json:"short_name"`
	Description string  `json:"description"`
	Weight      float64 `json:"weight"`
}

type EditGradeDataRequest struct {
	Number string `json:"number"`
	GradeDataRequest
}

func (g *GradeDataRequest) Validate() error {
	if g.Name == "" {
		return errors.New("name is required")
	}
	if g.Code == 0 {
		return errors.New("code is required")
	}

	return nil
}

func (g *EditGradeDataRequest) Validate() error {
	if g.Number == "" {
		return errors.New("number is required")
	}
	return g.GradeDataRequest.Validate()
}
