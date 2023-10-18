package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	subTopicModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
	"time"
)

type (
	SchemeOfWorkSubTopic struct {
		ID                 uint   `gorm:"primaryKey"`
		Number             string `json:"number" gorm:"size:36;unique;not null"`
		Code               int64  `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		EstimatedTime      int64  `json:"estimated_time"`
		SchemeOfWorkNumber string `json:"scheme_of_work_number" gorm:"size:36;not null"`
		SubtopicNumber     string `json:"sub_topic_number" gorm:"size:36;not null"`
		UserNumber         string `json:"user_number" gorm:"size:36;not null"`
		UpdatedAt          time.Time
		SchemeOfWork       SchemeOfWork            `gorm:"foreignKey:SchemeOfWorkNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		SubTopic           subTopicModels.SubTopic `gorm:"foreignKey:SubtopicNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User               models.User             `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
