package models

import (
	curriculumModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
	"gorm.io/gorm"
)
import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"

type CurriculumSubtopicReference struct {
	gorm.Model
	Number                 string                     `json:"number" gorm:"size:36;unique;not null;index"`
	Code                   int64                      `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
	LearningMaterialNumber string                     `json:"learning_material_number" gorm:"size:36;not null"`
	CurriculumNumber       string                     `json:"curriculum_number" gorm:"size:36;not null"`
	SubTopicNumber         string                     `json:"sub_topic_number" gorm:"size:36;not null"`
	UserNumber             string                     `json:"user_number" gorm:"size:36;not null"`
	LearningMaterial       LearningMaterial           `gorm:"foreignKey:LearningMaterialNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Curriculum             curriculumModel.Curriculum `gorm:"foreignKey:CurriculumNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SubTopic               models3.SubTopic           `gorm:"foreignKey:SubTopicNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	User                   models.User                `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

/*func (l *CurriculumSubtopicReference) Validate() error {
	//validate all the required fields here

	if l.Name == "" {
		return errors.New("name is required")
	}

	if l.Type == "" {
		return errors.New("type is required")
	}

	if l.UserNumber == "" {
		return errors.New("user number is required")
	}

	if l.ReferenceType == "" {
		return errors.New("reference type is required")
	}

	return nil
}*/
