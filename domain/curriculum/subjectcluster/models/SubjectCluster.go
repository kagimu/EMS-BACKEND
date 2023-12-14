package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	learningAreaModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/models"
	"time"
)

type (
	SubjectCluster struct {
		ID                 uint                            `gorm:"primaryKey"`
		Number             string                          `json:"number" gorm:"size:36;unique;not null"`
		Name               string                          `json:"name" form:"name" binding:"required" gorm:"type:longtext;not null"`
		Code               int64                           `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		CurriculumNumber   string                          `json:"curriculum_number" gorm:"size:36;not null"`
		UserNumber         string                          `json:"user_number" gorm:"size:40;not null"`
		LearningAreaNumber string                          `json:"learning_area_number" gorm:"size:40;not null"`
		CreatedAt          time.Time                       `json:"created_at"`
		UpdatedAt          time.Time                       `json:"updated_at"`
		DeletedAt          *time.Time                      `sql:"index"`
		Curriculum         models2.Curriculum              `gorm:"foreignKey:CurriculumNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		LearningArea       learningAreaModels.LearningArea `gorm:"foreignKey:LearningAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User               authModels.User                 `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *SubjectCluster) TableName() string {
	return "subject_clusters"
}
