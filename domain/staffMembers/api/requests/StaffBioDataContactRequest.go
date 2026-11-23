package requests

import (
	"errors"
)

type StaffBioDataContactRequest struct {
	MainEmail       string `json:"mainEmail"`
	SecondEmail     string `json:"secondEmail"`
	AltEmail        string `json:"altEmail"`
	MainPhone       string `json:"mainPhone"`
	SecondPhone     string `json:"secondPhone"`
	AltPhone        string `json:"altPhone"`
	PostMailAddress string `json:"postMailAddress"`
	PhysicalAddress string `json:"physicalAddress"`
}

func (r *StaffBioDataContactRequest) Validate() error {
	if r.MainEmail == "" {
		return errors.New("mainEmail is required")
	}
	if r.MainPhone == "" {
		return errors.New("mainPhone is required")
	}
	if r.PhysicalAddress == "" {
		return errors.New("physicalAddress is required")
	}
	return nil
}
