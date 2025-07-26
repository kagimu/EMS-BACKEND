package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	schemeOfWorkModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/models"
	"time"
)

type (
	LessonPlan struct {
		CreatedAt            time.Time
		UpdatedAt            time.Time
		Number               string                                 `json:"number" gorm:"size:36;primaryKey"`
		SchemeSubTopicNumber string                                 `json:"scheme_sub_topic_number" form:"name" binding:"required" gorm:"size:36;not null"`
		Status               string                                 `json:"status" gorm:"size:10;not null"`
		UserNumber           string                                 `json:"user_number" gorm:"size:36;not null"`
		SchemeOfWorkSubTopic schemeOfWorkModel.SchemeOfWorkSubTopic `gorm:"foreignKey:SchemeSubTopicNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User                 models.User                            `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
