package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/models"
	"time"
)

type (
	Subject struct {
		ID                   uint                   `gorm:"primaryKey"`
		Number               string                 `json:"number" gorm:"size:36;unique;not null"`
		Name                 string                 `json:"name" form:"name" binding:"required" gorm:"size:200;unique;not null"`
		Code                 string                 `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		SubjectClusterNumber string                 `json:"subject_cluster_number" gorm:"size:36;not null"`
		UserNumber           string                 `json:"user_number" gorm:"size:36;not null"`
		CreatedAt            time.Time              `json:"created_at"`
		UpdatedAt            time.Time              `json:"updated_at"`
		DeletedAt            *time.Time             `sql:"index"`
		SubjectCluster       models3.SubjectCluster `gorm:"foreignKey:SubjectClusterNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User                 models.User            `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *Subject) TableName() string {
	return "subjects"
}
