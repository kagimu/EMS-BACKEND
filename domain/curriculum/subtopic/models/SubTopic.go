package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	curriculumModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	topicModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/models"
	"time"
)

type (
	SubTopic struct {
		ID               uint                       `gorm:"primaryKey"`
		Number           string                     `json:"number" gorm:"size:36;unique;not null"`
		Name             string                     `json:"name" form:"name" binding:"required" gorm:"type:longtext;not null"`
		Code             int64                      `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		CurriculumNumber string                     `json:"curriculum_number" gorm:"size:36;not null"`
		UserNumber       string                     `json:"user_number" gorm:"size:36;not null"`
		EstimatedTime    int64                      `json:"estimated_time"`
		TopicNumber      string                     `json:"topic_number" gorm:"size:36;not null"`
		CreatedAt        time.Time                  `json:"created_at"`
		UpdatedAt        time.Time                  `json:"updated_at"`
		DeletedAt        *time.Time                 `sql:"index"`
		Curriculum       curriculumModel.Curriculum `gorm:"foreignKey:CurriculumNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Topic            topicModels.Topic          `gorm:"foreignKey:TopicNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User             models.User                `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *SubTopic) TableName() string {
	return "sub_topics"
}
