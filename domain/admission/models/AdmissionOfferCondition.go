package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type AdmissionOfferCondition struct {
	gorm.Model
	Number                   string          `json:"number" gorm:"size:36;unique;not null"`
	Description              string          `json:"description" gorm:"type:text;default:''"`
	UserNumber               string          `json:"user_number" gorm:"size:36;not null"`
	CustomAdmissionCondition string          `json:"custom_admission_condition" gorm:"type:text"`
	User                     authModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type AdmissionOfferConditions struct {
	gorm.Model
	Number                        string                  `json:"number" gorm:"size:36;unique;not null"`
	AdmissionOfferNumber          string                  `json:"admission_offer_number" gorm:"size:36"`
	AdmissionOffer                AdmissionOffer          `gorm:"foreignKey:AdmissionOfferNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AdmissionOfferConditionNumber string                  `json:"admission_offer_condition_number" gorm:"size:36"`
	AdmissionOfferCondition       AdmissionOfferCondition `gorm:"foreignKey:AdmissionOfferConditionNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber                    string                  `json:"user_number" gorm:"size:36;not null"`
	User                          authModels.User         `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
