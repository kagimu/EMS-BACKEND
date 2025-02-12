package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	studentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type Discount struct {
	gorm.Model
	Number         string                `json:"number" gorm:"size:36;unique;not null"`
	Reason         string                `json:"reason" gorm:"type:text"`
	DiscountAmount uint64                `json:"discount_amount"`
	AuthorizedBy   string                `json:"authorized_by" gorm:"size:36;not null"`
	Authorizor     models.User           `gorm:"foreignKey:AuthorizedBy;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudentNumber  string                `json:"student_number" gorm:"size:36;not null"`
	Student        studentModels.Student `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ServiceNumber  string                `json:"service_number" gorm:"size:36;not null"`
	Service        Service               `gorm:"foreignKey:ServiceNumber;references:ItemNumber;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber     string                `json:"user_number" gorm:"size:36;index;not null"`
	User           models.User           `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
