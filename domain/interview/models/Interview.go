package models

import (
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"gorm.io/gorm"
)

type Interview struct {
	gorm.Model
	Number            string                             `json:"number" gorm:"size:36;unique;not null"`
	Message           string                             `json:"message"`
	Venue             string                             `json:"venue" gorm:"size:100;not null"`
	Name              string                             `json:"name" gorm:"size:100;"`
	InterviewDate     string                             `json:"interview_date" gorm:"size:100;not null"`
	UserNumber        string                             `json:"user_number" gorm:"size:36;not null"`
	CampusClassNumber string                             `json:"campus_class_number" gorm:"size:36;not null"`
	CampusClass       models.CampusClass                 `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	User              authModel.User                     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AssessmentMetrics []assessmentModel.AssessmentMetric `gorm:"-"`
}

func (i *Interview) AfterFind(tx *gorm.DB) (err error) {
	// find assessment metrics for this interview
	var interviewAssessmentMetrics []InterviewAssessmentMetric

	if err := tx.Where("interview_number = ?", i.Number).Find(&interviewAssessmentMetrics).Error; err != nil {
		return err
	}

	for _, interviewAssessmentMetric := range interviewAssessmentMetrics {
		var assessmentMetric assessmentModel.AssessmentMetric
		if err := tx.First(&assessmentMetric, "number = ?", interviewAssessmentMetric.AssessmentMetricNumber).Error; err != nil {
			return err
		}
		i.AssessmentMetrics = append(i.AssessmentMetrics, assessmentMetric)
	}
	return nil
}
