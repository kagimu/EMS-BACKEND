package requests

import "errors"

type EntityRequest struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Facebook  string `json:"facebook"`
	Instagram string `json:"instagram"`
	Twitter   string `json:"twitter"`
	Website   string `json:"website"`
	Whatsapp  string `json:"whatsapp"`
	Linkedin  string `json:"linkedin"`
}

type UpdateEntityRequest struct {
	EntityRequest
	EntityNumber string `json:"entity_number"`
}

func (entityRequest *EntityRequest) Validate() error {

	if entityRequest.Name == "" {
		return errors.New("provide entity name")
	}

	if entityRequest.Phone == "" {
		return errors.New("provide entity phone")
	}

	if entityRequest.Email == "" {
		return errors.New("provide entity email")
	}

	if entityRequest.Address == "" {
		return errors.New("provide entity address")
	}

	return nil
}

func (entityRequest *UpdateEntityRequest) Validate() error {

	if err := entityRequest.EntityRequest.Validate(); err != nil {
		return err
	}

	if entityRequest.EntityNumber == "" {
		return errors.New("provide entity number being edited")
	}

	return nil
}
