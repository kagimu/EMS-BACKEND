package requests

import (
	"errors"
)

type StaffBioDataPersonalRequest struct {
	Surname            string `json:"surname"`
	SecondName         string `json:"secondName"`
	OtherNames         string `json:"otherNames"`
	DateCompleted      string `json:"dateCompleted"`
	DateOfBirth        string `json:"dateOfBirth"`
	Sex                string `json:"sex"`
	BirthCertificateNo string `json:"birthCertificateNo"`
	Religion           string `json:"religion"`
	NationalIdNo       string `json:"nationalIdNo"`
	CountryOfBirth     string `json:"countryOfBirth"`
	HomeCountry        string `json:"homeCountry"`
	HomeDistrict       string `json:"homeDistrict"`
	HomeCounty         string `json:"homeCounty"`
	HomeSubCounty      string `json:"homeSubCounty"`
	HomeParish         string `json:"homeParish"`
	HomeVillage        string `json:"homeVillage"`
	VillageOfResidence string `json:"villageOfResidence"`
	TaxIdNo            string `json:"taxIdNo"`
}

func (r *StaffBioDataPersonalRequest) Validate() error {
	if r.Surname == "" {
		return errors.New("surname is required")
	}
	if r.SecondName == "" {
		return errors.New("secondName is required")
	}
	if r.DateOfBirth == "" {
		return errors.New("dateOfBirth is required")
	}
	if r.Sex == "" {
		return errors.New("sex is required")
	}
	if r.NationalIdNo == "" {
		return errors.New("nationalIdNo is required")
	}
	if r.HomeDistrict == "" {
		return errors.New("homeDistrict is required")
	}
	return nil
}
