package models

import (
	learningMaterialModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/learningMaterials/models"
	"time"
)

type LessonPlanReference struct {
	CreatedAt              time.Time
	UpdatedAt              time.Time
	Number                 string                                 `json:"number" gorm:"size:36;primaryKey"`
	Code                   int64                                  `json:"code" form:"code" binding:"required" gorm:"size:80;not null"`
	Name                   string                                 `json:"name" form:"name" binding:"required" gorm:"size:80;not null"`
	LessonPlanNumber       string                                 `json:"lesson_plan_number" gorm:"size:36;not null"`
	LearningMaterialNumber string                                 `json:"learning_material_number" gorm:"size:36;not null"`
	UserNumber             string                                 `json:"user_number" gorm:"size:36;not null"`
	LessonPlan             LessonPlan                             `gorm:"foreignKey:LessonPlanNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	LearningMaterial       learningMaterialModel.LearningMaterial `gorm:"foreignKey:LearningMaterialNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
