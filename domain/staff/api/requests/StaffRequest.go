package requests

import "errors"

type StaffRequest struct {
	StaffNumber                string `json:"staff_number"`
	Surname                    string `json:"surname"`
	FirstName                  string `json:"first_name"`
	OtherName                  string `json:"other_name"`
	Gender                     string `json:"gender"`
	MainTelephoneNumber        string `json:"main_telephone_number"`
	AlternativeTelephoneNumber string `json:"alternative_telephone_number"`
	MainEmailAddress           string `json:"main_email_address"`
	AlternativeEmailAddress    string `json:"alternative_email_address"`
	JobTitleNumber             string `json:"job_title_number"`
	CampusNumber               string `json:"campus_number"`
	SalaryScaleNumber          string `json:"salary_scale_number"`
}

func (s *StaffRequest) Validate() error {
	if s.Surname == "" {
		return errors.New("surname is required")
	}
	if s.FirstName == "" {
		return errors.New("first name is required")
	}
	if s.Gender == "" {
		return errors.New("gender is required")
	}
	if s.MainTelephoneNumber == "" {
		return errors.New("main telephone number is required")
	}
	if s.JobTitleNumber == "" {
		return errors.New("job title number is required")
	}
	if s.MainEmailAddress == "" {
		return errors.New("main email address is required")
	}
	if s.CampusNumber == "" {
		return errors.New("campus number is required")
	}

	return nil
}
