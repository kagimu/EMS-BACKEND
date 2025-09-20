package models

import (
	"time"

	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
)

type Narrative struct {
	Number                  string `gorm:"primarykey;size:36; not null;unique"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
	ReferenceDate           time.Time
	ReferenceTime           string                       `gorm:"size:8"`
	Description             string                       `gorm:"size:2500;not null"`
	Status                  string                       `gorm:"size:100"`
	CurrentStatus           string                       `gorm:"size:100"`
	StudentNumber           string                       `gorm:"size:36;index;not null;column:student_number"`
	Student                 studentModel.Student         `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	StudyPeriodNumber       string                       `gorm:"size:36;index;not null;column:study_period_number"`
	StudyPeriod             studyPeriodModel.StudyPeriod `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserNumber              string                       `gorm:"size:36;index;not null;column:user_number"`
	User                    authModel.User               `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusClassNumber       string                       `gorm:"size:36; not null;column:campus_class_number"`
	CampusClass             classModels.CampusClass      `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	AssessmentAreaNumber    string                       `gorm:"size:36; not null;column:assessment_area_number"`
	AssessmentArea          models.AssessmentArea        `gorm:"foreignKey:AssessmentAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	SuggestedSolution       string                       `gorm:"size:2500;not null"`
	NarrativeScheduleNumber string                       `gorm:"size:36;not null"`
}

type ScheduledNarrative struct {
	Number               string `gorm:"primarykey;not null;size:36;unique"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	ScheduledOn          time.Time
	Description          string                       `gorm:"size:2500"`
	StudyPeriodNumber    string                       `gorm:"size:36;index;not null;column:study_period_number"`
	StudyPeriod          studyPeriodModel.StudyPeriod `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserNumber           string                       `gorm:"size:36;index;not null;column:user_number"`
	User                 authModel.User               `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusClassNumber    string                       `gorm:"size:36; not null;column:campus_class_number"`
	CampusClass          classModels.CampusClass      `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	AssessmentAreaNumber string                       `gorm:"size:36; not null;column:assessment_area_number"`
	AssessmentArea       models.AssessmentArea        `gorm:"foreignKey:AssessmentAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

type NarrativeSchedule struct {
	ScheduledNarrativeNumber string `gorm:"size:36;primarykey;column:scheduled_narrative_number"`
	NarrativeNumber          string `gorm:"size:36;primarykey;column:narrative_number"`

	ScheduledNarrative ScheduledNarrative `gorm:"foreignKey:ScheduledNarrativeNumber;references:Number;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Narrative          Narrative          `gorm:"foreignKey:NarrativeNumber;references:Number;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ScheduledNarrativeScriptMarkingGuide struct {
	CreatedAt                time.Time
	UpdatedAt                time.Time
	ScriptNumber             string             `gorm:"column:script_number;size:36;unique"`
	Script                   models2.Attachment `gorm:"foreignKey:ScriptNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	MarkingGuideNumber       string             `gorm:"column:marking_guide_number;size:36;unique"`
	MarkingGuide             models2.Attachment `gorm:"foreignKey:MarkingGuideNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserNumber               string             `gorm:"size:36;not null;column:user_number"`
	User                     authModel.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ScheduledNarrativeNumber string             `gorm:"primaryKey;size:36;not null;column:scheduled_narrative_number"`
	ScheduledNarrative       ScheduledNarrative `gorm:"foreignKey:ScheduledNarrativeNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
