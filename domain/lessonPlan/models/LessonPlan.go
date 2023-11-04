package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	conceptModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	schemeOfWorkModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/schemeOfWork/models"
	"gorm.io/gorm"
)

type (
	LessonPlan struct {
		gorm.Model
		Number                  string                                 `json:"number" gorm:"size:36;unique;not null"`
		LessonPlanConceptNumber string                                 `json:"lesson_plan_concept_number" form:"name" binding:"required" gorm:"size:36;not null"`
		SchemeSubTopicNumber    string                                 `json:"scheme_sub_topic_number" form:"name" binding:"required" gorm:"size:36;not null"`
		Code                    int64                                  `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		UserNumber              string                                 `json:"user_number" gorm:"size:36;not null"`
		SchemeOfWork            schemeOfWorkModel.SchemeOfWorkSubTopic `gorm:"foreignKey:SchemeSubTopicNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Concept                 conceptModel.Concept                   `gorm:"foreignKey:LessonPlanConceptNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User                    models.User                            `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
