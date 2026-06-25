package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	schoolModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	departmentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/departments/models"
	"gorm.io/gorm"
)

type JobPosition struct {
	gorm.Model
	Number           string                      `json:"number" gorm:"size:36;unique;index;not null"`
	PostName         string                      `json:"post_name" form:"post_name" binding:"required" gorm:"size:100;not null"`
	PostCode         string                      `json:"post_code" form:"post_code" binding:"required" gorm:"size:50;not null"`
	JobCode          string                      `json:"job_code" gorm:"size:50;unique;index"`
	SchoolNumber     string                      `json:"school_number" gorm:"size:36;not null"`
	School           schoolModels.School         `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	DepartmentNumber string                      `json:"department_number" gorm:"size:36;not null"`
	Department       departmentModels.Department `gorm:"foreignKey:DepartmentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber       string                      `json:"user_number" gorm:"size:36;not null"`
	User             authModels.User             `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
