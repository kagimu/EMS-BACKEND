package models

import (
	authModels "github.com/hamsoftug/wholesomeEd/domain/auth/models"
	learningAreaModels "github.com/hamsoftug/wholesomeEd/domain/curriculum/learningarea/models"
	"time"
)

type (
	SubjectCluster struct {
		ID                 uint                            `gorm:"primaryKey"`
		Number             string                          `json:"number" gorm:"size:100;unique;not null"`
		Name               string                          `json:"name" form:"name" binding:"required" gorm:"size:60;unique;not null"`
		UserNumber         string                          `json:"user_number" gorm:"size:40;not null"`
		LearningAreaNumber string                          `json:"learning_area_number" gorm:"size:40;not null"`
		CreatedAt          time.Time                       `json:"created_at"`
		UpdatedAt          time.Time                       `json:"updated_at"`
		DeletedAt          *time.Time                      `sql:"index"`
		LearningArea       learningAreaModels.LearningArea `gorm:"foreignKey:LearningAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User               authModels.User                 `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *SubjectCluster) TableName() string {
	return "subject_clusters"
}
