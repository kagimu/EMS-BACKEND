package models

import (
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	curriculumModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
)

type (
	Concept struct {
		ID                uint                       `gorm:"primaryKey"`
		Number            string                     `json:"number" gorm:"size:36;unique;not null"`
		Name              string                     `json:"name" form:"name" binding:"required" gorm:"type:longtext;not null"`
		Code              int64                      `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		CurriculumNumber  string                     `json:"curriculum_number" gorm:"size:36;not null"`
		UserNumber        string                     `json:"user_number" gorm:"size:36;not null"`
		SubTopicNumber    string                     `json:"sub_topic_number" gorm:"size:36;not null"`
		EstimatedTime     int64                      `json:"estimated_time" gorm:"size:60"`
		RequiredMaterials string                     `json:"required_materials" gorm:"type:longtext"`
		CreatedAt         time.Time                  `json:"created_at"`
		UpdatedAt         time.Time                  `json:"updated_at"`
		Curriculum        curriculumModel.Curriculum `gorm:"foreignKey:CurriculumNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
		SubTopic          models3.SubTopic           `gorm:"foreignKey:SubTopicNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
		User              models.User                `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *Concept) TableName() string {
	return "concepts"
}
