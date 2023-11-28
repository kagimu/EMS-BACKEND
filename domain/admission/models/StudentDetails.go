package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type (
	StudentDetail struct {
		gorm.Model
		Number             string      `json:"number" gorm:"size:36;unique;not null"`
		Surname            string      `json:"surname" form:"surname" binding:"required" gorm:"size:80;not null"`
		SecondName         string      `json:"second_name" form:"second_name" binding:"required" gorm:"size:80;not null"`
		OtherName          string      `json:"other_name" form:"other_name" binding:"required" gorm:"size:80;not null"`
		DateOfBirth        string      `json:"date_of_birth" form:"date_of_birth" binding:"required" gorm:"not null"`
		BirthCertificateNo string      `json:"birth_certificate_no" form:"birth_certificate_no" binding:"required" gorm:"size:100;not null"`
		CountryOfBirth     string      `json:"country_of_birth" form:"country_of_birth" binding:"required" gorm:"size:100;not null"`
		HomeDistrict       string      `json:"home_district" form:"home_district" binding:"required" gorm:"size:100;not null"`
		HomeSubCounty      string      `json:"home_sub_county" form:"home_sub_county" binding:"required" gorm:"size:100;not null"`
		HomeVillage        string      `json:"home_village" form:"home_village" binding:"required" gorm:"size:100;not null"`
		HomePhone1         string      `json:"home_phone1" form:"home_phone1" binding:"required" gorm:"size:60;unique;not null"`
		HomePhone2         string      `json:"home_phone2" form:"home_phone2" binding:"required" gorm:"size:60;unique;not null"`
		HomeEmail1         string      `json:"home_email1" form:"home_email1" binding:"required" gorm:"size:60;not null"`
		HomeEmail2         string      `json:"home_email2" form:"home_email2" binding:"required" gorm:"size:60;not null"`
		Religion           string      `json:"religion" form:"religion" binding:"required" agorm:"size:60;not null"`
		HomeCountry        string      `json:"home_country" form:"home_country" binding:"required" gorm:"size:60;not null"`
		HomeCounty         string      `json:"home_county" form:"home_county" binding:"required" gorm:"size:60;not null"`
		HomeParish         string      `json:"home_parish" form:"home_parish" binding:"required" gorm:"size:60;not null"`
		NearestTown        string      `json:"nearest_town" form:"nearest_town" binding:"required" gorm:"size:60;not null"`
		UserNumber         string      `json:"user_number" gorm:"size:36;not null"`
		User               models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
