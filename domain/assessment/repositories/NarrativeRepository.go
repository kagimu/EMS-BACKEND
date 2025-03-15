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

func SaveNarrativeAndSchedule(narrativeSchedule models.NarrativeSchedule, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	return db.Save(&narrativeSchedule).Error
}

func FindNarrativeByScheduleNarrativeNumberAndStudentNumber(scheduleNarrativeNumber, studentNumber string, tx ...*gorm.DB) (models.Narrative, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var narrative models.Narrative

	if err := db.
		Joins("narrative_schedules on narratives.number = narrative_schedules.narrative_number").
		Where("narrative_schedules.schedule_narrative_number = ? AND narratives.student_number = ?", scheduleNarrativeNumber, studentNumber).Find(&narrative).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return narrative, nil
		}

		return narrative, err
	}

	return narrative, nil
}

func FindNarrativeByNumber(narrativeNumber string, tx ...*gorm.DB) (models.Narrative, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var narrative models.Narrative

	if err := db.Where("number = ?", narrativeNumber).
		Find(&narrative).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Narrative{}, nil
		}
		return models.Narrative{}, err
	}

	return narrative, nil
}

func FindScheduledNarrativeByNumber(scheduledNarrativeNumber string, tx ...*gorm.DB) (models.ScheduledNarrative, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var scheduledNarrative models.ScheduledNarrative

	if err := db.Where("number = ?", scheduledNarrativeNumber).
		Find(&scheduledNarrative).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ScheduledNarrative{}, nil
		}
		return models.ScheduledNarrative{}, err
	}

	return scheduledNarrative, nil
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

func FindAllScheduledNarratives(studyPeriodNumber, campusClassNumber, assessmentAreaNumber string, startDate, endDate time.Time, tx ...*gorm.DB) ([]models.ScheduledNarrative, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var narratives []models.ScheduledNarrative

	if err := db.
		Scopes(filters.ByAnyField("study_period_number", studyPeriodNumber), filters.ByDateBetween("scheduled_on", startDate, endDate), filters.ByAnyField("assessment_area_number", assessmentAreaNumber), filters.ByAnyField("campus_class_number", campusClassNumber)).
		Find(&narratives).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.ScheduledNarrative{}, nil
		}
		return nil, err
	}

	return narratives, nil
}
