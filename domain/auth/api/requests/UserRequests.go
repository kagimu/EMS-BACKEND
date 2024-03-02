package requests

import "errors"

type UserDataRequest struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	LastName   string `json:"last_name"`
	Gender     string `json:"gender"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
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

	return nil
}
