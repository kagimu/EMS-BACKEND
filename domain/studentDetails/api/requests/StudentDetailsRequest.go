package requests

import "errors"

type StudentDetailsRequest struct {
	Surname            string `json:"surname"`
	SecondName         string `json:"second_name"`
	OtherName          string `json:"other_name"`
	DateOfBirth        string `json:"date_of_birth"`
	BirthCertificateNo string `json:"birth_certificate_no"`
	CountryOfBirth     string `json:"country_of_birth"`
	HomeDistrict       string `json:"home_district"`
	HomeSubCounty      string `json:"home_sub_county"`
	HomeVillage        string `json:"home_village"`
	HomePhone1         string `json:"home_phone1"`
	HomeEmail1         string `json:"home_email1"`
	Religion           string `json:"religion"`
	HomeCountry        string `json:"home_country"`
	HomeCounty         string `json:"home_county"`
	HomeParish         string `json:"home_parish"`
	NearestTown        string `json:"nearest_town"`
	HomePhone2         string `json:"home_phone2"`
	HomeEmail2         string `json:"home_email2"`
}

func (s *StudentDetailsRequest) Validate() error {
	if s.Surname == "" {
		return errors.New("surname is required")
	}
	if s.SecondName == "" {
		return errors.New("second name is required")
	}
	if s.DateOfBirth == "" {
		return errors.New("date of birth is required")
	}
	if s.BirthCertificateNo == "" {
		return errors.New("birth certificate no is required")
	}
	if s.CountryOfBirth == "" {
		return errors.New("country of birth is required")
	}
	if s.HomeDistrict == "" {
		return errors.New("home district is required")
	}
	if s.HomeSubCounty == "" {
		return errors.New("home sub county is required")
	}
	if s.HomeVillage == "" {
		return errors.New("home village is required")
	}
	if s.HomeCounty == "" {
		return errors.New("home county is required")
	}
	if s.HomePhone1 == "" {
		return errors.New("hhme phone 1 is required")
	}
	if s.HomeEmail1 == "" {
		return errors.New("home email 1 is required")
	}
	if s.Religion == "" {
		return errors.New("religion is required")
	}
	if s.HomeCountry == "" {
		return errors.New("home country is required")
	}
	if s.HomeParish == "" {
		return errors.New("home parish is required")
	}
	if s.NearestTown == "" {
		return errors.New("nearest town is required")
	}

	return nil
}
