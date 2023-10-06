package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"time"
)

type (
	ImportCurriculum struct {
		ID         uint        `gorm:"primaryKey"`
		Number     string      `json:"number" gorm:"size:36;unique;not null"`
		Name       string      `json:"name" form:"name" binding:"required" gorm:"size:200;not null"`
		FileName   string      `json:"file_name" form:"file_name" binding:"required" gorm:"size:200;not null"`
		UserNumber string      `json:"user_number" gorm:"size:36;not null"`
		CreatedAt  time.Time   `json:"created_at"`
		UpdatedAt  time.Time   `json:"updated_at"`
		DeletedAt  *time.Time  `sql:"index"`
		User       models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}

	ExcelCurriculumInfo struct {
		LearningArea       string
		SubjectCluster     string
		Subject            string
		Class              string
		Strand             string
		Topic              string
		SubTopic           string
		Concept            string
		LearningAreaCode   string
		SubjectClusterCode string
		SubjectCode        string
		StrandCode         string
		TopicCode          string
		SubTopicCode       string
		EstimatedTime      string
		ConceptCode        string
	}

	/*SubjectsInfo struct {
		Number string `json:"number"`
		Name   string `json:"name"`
	}*/
)

func (u *ImportCurriculum) TableName() string {
	return "import_curriculums"
}
