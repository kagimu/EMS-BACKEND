package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
	"time"
)

func SaveNarrative(narrative models.Narrative, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	return db.Save(&narrative).Error
}

func SaveScheduledNarrative(scheduledNarrative models.ScheduledNarrative, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	return db.Save(&scheduledNarrative).Error
}

func FindAllNarratives(studentNumber, studyPeriodNumber, campusClassNumber, assessmentAreaNumber string, startDate, endDate time.Time, tx ...*gorm.DB) ([]models.Narrative, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var narratives []models.Narrative

	if err := db.Where("student_number = ?", studentNumber).
		Scopes(filters.ByAnyField("study_period_number", studyPeriodNumber), filters.ByDateBetween("created_at", startDate, endDate), filters.ByAnyField("assessment_area_number", assessmentAreaNumber), filters.ByAnyField("campus_class_number", campusClassNumber)).
		Find(&narratives).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Narrative{}, nil
		}
		return nil, err
	}

	return narratives, nil
}
func FindAllScheduledNarratives(studyPeriodNumber, campusClassNumber, assessmentAreaNumber string, startDate time.Time, tx ...*gorm.DB) ([]models.ScheduledNarrative, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var narratives []models.ScheduledNarrative

	if err := db.
		Scopes(filters.ByAnyField("study_period_number", studyPeriodNumber), filters.ByDateAfter("scheduled_on", startDate), filters.ByAnyField("assessment_area_number", assessmentAreaNumber), filters.ByAnyField("campus_class_number", campusClassNumber)).
		Find(&narratives).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.ScheduledNarrative{}, nil
		}
		return nil, err
	}

	return narratives, nil
}
