package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type (
	StudentDetail struct {
		gorm.Model
		Number string `json:"number" gorm:"size:36;unique;not null"`

		HomeVillage string `json:"home_village"`
		HomePhone1  string `json:"home_phone1"`
		HomeEmail1  string `json:"home_email1"`
		Religion    string `json:"religion"`
		HomeCountry string `json:"home_country"`
		HomeCounty  string `json:"home_county"`
		HomeParish  string `json:"home_parish"`
		NearestTown string `json:"nearest_town"`
		HomePhone2  string `json:"home_phone2"`
		HomeEmail2  string `json:"home_email2"`

		Surname            string `json:"surname" form:"surname" binding:"required" gorm:"size:80;not null"`
		SecondName         string `json:"second_name" form:"second_name" binding:"required" gorm:"size:80;not null"`
		OtherName          string `json:"other_name" form:"other_name" binding:"required" gorm:"size:80;not null"`
		DateOfBirth        string `json:"date_of_birth" form:"date_of_birth" binding:"required" gorm:"not null"`
		BirthCertificateNo string `json:"birth_certificate_no" form:"birth_certificate_no" binding:"required" gorm:"size:100;not null"`
		CountryOfBirth     string `json:"country_of_birth" form:"country_of_birth" binding:"required" gorm:"size:100;not null"`
		HomeDistrict       string `json:"home_district" form:"home_district" binding:"required" gorm:"size:100;not null"`
		HomeSubCounty      string `json:"home_sub_county" form:"home_sub_county" binding:"required" gorm:"size:100;not null"`

		/*Phone        string      `json:"phone" form:"phone" binding:"required" gorm:"size:40;unique;not null"`
		Email        string      `json:"email" form:"email" binding:"required" gorm:"size:40;not null"`
		Location     string      `json:"location" form:"location" binding:"required" gorm:"size:200;not null"`
		UserNumber   string      `json:"user_number" gorm:"size:36;not null"`
		SchoolNumber string      `json:"school_number" gorm:"size:36;not null"`*/
		User models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
