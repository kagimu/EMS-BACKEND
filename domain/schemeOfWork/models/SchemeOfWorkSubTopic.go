package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	subTopicModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
	"gorm.io/gorm"
)

type (
	SchemeOfWorkSubTopic struct {
		gorm.Model
		ID                 uint                    `gorm:"primaryKey"`
		Number             string                  `json:"number" gorm:"size:36;unique;not null"`
		Code               string                  `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		SchemeOfWorkNumber string                  `json:"scheme_of_work_number" gorm:"size:36;not null"`
		SubtopicNumber     string                  `json:"sub_topic_number" gorm:"size:36;not null"`
		UserNumber         string                  `json:"user_number" gorm:"size:36;not null"`
		SchemeOfWork       SchemeOfWork            `gorm:"foreignKey:SchemeOfWorkNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		SubTopic           subTopicModels.SubTopic `gorm:"foreignKey:SubtopicNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User               models.User             `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
