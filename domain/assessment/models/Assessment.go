package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	assessmentAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	levelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type Assessment struct {
	gorm.Model
	Number                string                             `json:"number" gorm:"size:36;unique;not null"`
	Subject               string                             `json:"subject" gorm:"type:text"`
	Attendance            uint64                             `json:"attendance"`
	Observations          string                             `json:"observations" gorm:"type:text"`
	GuidedExercise        uint64                             `json:"guided_exercise"`
	UnguidedExercise      uint64                             `json:"unguided_exercise"`
	ExplorativeAssignment uint64                             `json:"explorative_assignment"`
	GroupWork             uint64                             `json:"group_work"`
	TopicalText           uint64                             `json:"topical_text"`
	Project               uint64                             `json:"project"`
	CreativeExperiment    uint64                             `json:"creative_experiment"`
	AssessmentAreaNumber  string                             `json:"assessment_area_number"`
	StudentNumber         string                             `json:"student_number" gorm:"size:36;index;not null"`
	Student               studentModel.Student               `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AssessmentArea        assessmentAreaModel.AssessmentArea `gorm:"foreignKey:AssessmentAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	LevelNumber           string                             `json:"level_number" gorm:"size:36;not null"`
	Level                 levelModels.Level                  `gorm:"foreignKey:LevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ExaminationTypeNumber string                             `json:"examination_type_number" gorm:"size:36;"`
	Examination           ExaminationType                    `gorm:"foreignKey:ExaminationTypeNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber            string                             `json:"user_number" gorm:"size:36;index;not null"`
	User                  authModel.User                     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type ExaminationType struct {
	gorm.Model
	Number string `json:"number" gorm:"size:36;unique;not null"`
	Name   string `json:"name" gorm:"type:text;size:20"`
}
