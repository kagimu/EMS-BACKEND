package requests

import "errors"

type EditRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

type NewRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

func (req *NewRoleRequest) Validate() error {
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
