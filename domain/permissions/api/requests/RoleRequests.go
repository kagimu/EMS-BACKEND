package requests

import "errors"

type EditRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

type NewRoleRequest struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Permissions  []string `json:"permissions"`
	EntityNumber string   `json:"entity_number"`
}

func (req *NewRoleRequest) Validate() error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if len(req.Name) < 3 {
		return errors.New("the name of the role must be at least 3 characters")
	}

	if len(req.Name) > 50 {
		return errors.New("the name of the role must not exceed 50 characters")
	}

	if req.Description == "" {
		return errors.New("description is required")
	}
	if len(req.Permissions) == 0 {
		return errors.New("at least one permission is required")
	}
	if req.EntityNumber == "" {
		return errors.New("entity number is required")
	}

	return nil
}

func (req *EditRoleRequest) Validate() error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Description == "" {
		return errors.New("description is required")
	}
	if len(req.Permissions) == 0 {
		return errors.New("at least one permission is required")
	}
	return nil

}
