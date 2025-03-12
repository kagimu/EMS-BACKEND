package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"time"
)

type Narrative struct {
	Number               string `gorm:"primarykey;size:36; not null;unique"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	ReferenceDate        time.Time
	ReferenceTime        string                       `gorm:"size:8"`
	Description          string                       `gorm:"size:2500;not null"`
	StudentNumber        string                       `gorm:"size:36;index;not null;column:student_number"`
	Student              studentModel.Student         `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	StudyPeriodNumber    string                       `gorm:"size:36;index;not null;column:study_period_number"`
	StudyPeriod          studyPeriodModel.StudyPeriod `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserNumber           string                       `gorm:"size:36;index;not null;column:user_number"`
	User                 authModel.User               `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusClassNumber    string                       `gorm:"size:36; not null;column:campus_class_number"`
	CampusClass          classModels.CampusClass      `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	AssessmentAreaNumber string                       `gorm:"size:36; not null;column:assessment_area_number"`
	AssessmentArea       models.AssessmentArea        `gorm:"foreignKey:AssessmentAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

type ScheduledNarrative struct {
	Number               string `gorm:"primarykey;not null;size:36;unique"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	ScheduledOn          time.Time
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
