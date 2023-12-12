package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type (
	CampusClass struct {
		gorm.Model
		Number                string              `json:"number" gorm:"size:36;unique;not null"`
		CurriculumClassNumber string              `json:"curriculum_class_number" gorm:"size:36;not null"`
		CampusNumber          string              `json:"campus_number" gorm:"size:36;not null"`
		CampusLevelNumber     string              `json:"campus_level_number" gorm:"size:36;not null"`
		Code                  int64               `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		UserNumber            string              `json:"user_number" gorm:"size:36;not null"`
		CurriculumClass       CurriculumClass     `gorm:"foreignKey:CurriculumClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Campus                models2.Campus      `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		CampusLevel           models3.CampusLevel `gorm:"foreignKey:CampusLevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User                  models.User         `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		ClassName             string              `json:"class_name" gorm:"-"`
		LevelName             string              `json:"level_name" gorm:"-"`
		Streams               []Stream            `gorm:"foreignKey:ClassLevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (classLevel *CampusClass) AfterFind(tx *gorm.DB) (err error) {
	var class Class
	//tx.Model(&classLevel.Class).Where("number = ?", classLevel.ClassNumber).First(&class)
	classLevel.ClassName = class.Name

	return
}

func (classLevel *CampusClass) PreInsert(tx *gorm.DB) (err error) {
	classLevel.Number = uuid.NewV4().String()

	return
}
