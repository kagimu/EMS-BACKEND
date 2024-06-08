package requests

import (
	"errors"
	"html"
	"strings"
)

type UserDataRequest struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	LastName   string `json:"last_name"`
	Gender     string `json:"gender"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}

func (ur *UserDataRequest) Prepare() {
	ur.FirstName = html.EscapeString(strings.TrimSpace(ur.FirstName))
	ur.SecondName = html.EscapeString(strings.TrimSpace(ur.SecondName))
	ur.LastName = html.EscapeString(strings.TrimSpace(ur.LastName))
	ur.Email = html.EscapeString(strings.TrimSpace(ur.Email))
	ur.Phone = html.EscapeString(strings.TrimSpace(ur.Phone))
}

func (ur *UserDataRequest) Validate() error {

	if ur.FirstName == "" {
		return errors.New("first Name is required")
	}

	if ur.SecondName == "" {
		return errors.New("second Name is required")
	}

	if ur.Gender == "" {
		return errors.New("gender is required")
	}

	if ur.Email == "" {
		return errors.New("email is required")
	}

	if ur.Phone == "" {
		return errors.New("phone is required")
	}

	validMale := ur.Gender != "Male"
	validFemale := ur.Gender != "Female"

	if validMale && validFemale {
		return errors.New("invalid gender")
	}

	return nil
}

type StaffAccountRequest struct {
	FirstName       string   `json:"first_name"`
	SecondName      string   `json:"second_name"`
	LastName        string   `json:"last_name"`
	Gender          string   `json:"gender"`
	Email           string   `json:"email"`
	Phone           string   `json:"phone"`
	RoleNumbers     []string `json:"role_numbers"`
	Password        string   `json:"password"`
	ConfirmPassword string   `json:"confirm_password"`
	EntityNumber    string   `json:"entity_number"`
}

func (sr *StaffAccountRequest) ValidateStaffAccountRequest() error {
	if sr.FirstName == "" {
		return errors.New("first name is required")
	}
	if sr.SecondName == "" {
		return errors.New("second name is required")
	}
	if sr.LastName == "" {
		return errors.New("last name is required")
	}
	if sr.Gender == "" {
		return errors.New("gender is required")
	}
	if sr.Email == "" {
		return errors.New("email is required")
	}
	if sr.Phone == "" {
		return errors.New("phone is required")
	}
	if len(sr.RoleNumbers) < 1 {
		return errors.New("please select at least one role")
	}

	if sr.Password == "" {
		return errors.New("password is required")
	}
	if sr.ConfirmPassword == "" {
		return errors.New("confirm password is required")
	}
	if sr.Password != sr.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	if sr.EntityNumber == "" {
		return errors.New("entity number is required")
	}

	return nil
}
