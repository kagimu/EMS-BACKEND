package requests

import (
	"errors"
	"html"
	"strings"
)

type UserPassword struct {
	Token           string `json:"token"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (ur *UserPassword) Validate() error {

	if ur.Token == "" {
		return errors.New("token is required")
	}

	if ur.Password == "" {
		return errors.New("password is required")
	}

	if ur.Password != ur.ConfirmPassword {
		return errors.New("password Mismatch")
	}

	return nil
}

type UserEmail struct {
	Email string `json:"email"`
}

func (ur *UserEmail) Prepare() {
	ur.Email = html.EscapeString(strings.TrimSpace(ur.Email))
}

func (ur *UserEmail) Validate() error {

	if ur.Email == "" {
		return errors.New("email is required")
	}

	return nil
}
