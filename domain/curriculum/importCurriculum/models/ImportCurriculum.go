package models

import (
	"html"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type (
	ImportCurriculum struct {
		gorm.Model
		Number           string      `json:"number" gorm:"size:36;unique;not null"`
		Name             string      `json:"name" form:"name" binding:"required" gorm:"size:200;not null"`
		FileName         string      `json:"file_name" form:"file_name" binding:"required" gorm:"size:200;not null"`
		CurriculumNumber string      `json:"curriculum_number" gorm:"size:36;not null"`
		UserNumber       string      `json:"user_number" gorm:"size:36;not null"`
		Curriculum       Curriculum  `gorm:"foreignKey:CurriculumNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User             models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}

	ExcelCurriculumInfo struct {
		LearningArea         string
		SubjectCluster       string
		Subject              string
		Class                string
		Strand               string
		Topic                string
		SubTopic             string
		Concept              string
		LearningAreaCode     string
		SubjectClusterCode   string
		SubjectCode          string
		StrandCode           string
		TopicCode            string
		SubTopicCode         string
		EstimatedTime        string
		ConceptCode          string
		TeachingSequenceCode string
	}
)

func (p *ExcelCurriculumInfo) Prepare() {
	p.LearningArea = html.EscapeString(strings.TrimSpace(p.LearningArea))
	p.SubjectCluster = html.EscapeString(strings.TrimSpace(p.SubjectCluster))
	p.Subject = html.EscapeString(strings.TrimSpace(p.Subject))
	p.Class = html.EscapeString(strings.TrimSpace(p.Class))
	p.Strand = html.EscapeString(strings.TrimSpace(p.Strand))
	p.Topic = html.EscapeString(strings.TrimSpace(p.Topic))
	p.SubTopic = html.EscapeString(strings.TrimSpace(p.SubTopic))
	p.Concept = html.EscapeString(strings.TrimSpace(p.Concept))
	p.LearningAreaCode = html.EscapeString(strings.TrimSpace(p.LearningAreaCode))
	p.SubjectClusterCode = html.EscapeString(strings.TrimSpace(p.SubjectClusterCode))
	p.SubjectCode = html.EscapeString(strings.TrimSpace(p.SubjectCode))
	p.StrandCode = html.EscapeString(strings.TrimSpace(p.StrandCode))
	p.TopicCode = html.EscapeString(strings.TrimSpace(p.TopicCode))
	p.SubTopicCode = html.EscapeString(strings.TrimSpace(p.SubTopicCode))
	p.EstimatedTime = html.EscapeString(strings.TrimSpace(p.EstimatedTime))
	p.ConceptCode = html.EscapeString(strings.TrimSpace(p.ConceptCode))
	p.TeachingSequenceCode = html.EscapeString(strings.TrimSpace(p.TeachingSequenceCode))
}

func (u *ImportCurriculum) TableName() string {
	return "import_curriculums"
}
