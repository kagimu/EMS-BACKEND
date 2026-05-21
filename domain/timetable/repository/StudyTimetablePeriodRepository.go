package repository

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"gorm.io/gorm"
)

type StudyTimetablePeriodRepository interface {
	SaveStudyTimetablePeriod(period models.StudyTimetablePeriod, txs ...*gorm.DB) error
	UpdateStudyTimetablePeriod(period *models.StudyTimetablePeriod) error
	DeleteStudyTimetablePeriod(period models.StudyTimetablePeriod) error
	FindStudyTimetablePeriodByNumber(number string) (models.StudyTimetablePeriod, error)
	FindStudyTimetablePeriods(studyTimetableNumber, streamNumber, day string) ([]models.StudyTimetablePeriod, error)
	ExistsStudyTimetablePeriodByDayAndTime(streamNumber, studyPeriodNumber, day, startTime, endTime string) (bool, error)
}

type studyTimetablePeriodRepository struct {
	DB *gorm.DB
}

func NewStudyTimetablePeriodRepository(db *gorm.DB) StudyTimetablePeriodRepository {
	return &studyTimetablePeriodRepository{DB: db}
}

func (r *studyTimetablePeriodRepository) SaveStudyTimetablePeriod(period models.StudyTimetablePeriod, txs ...*gorm.DB) error {
	db := r.DB
	if len(txs) > 0 {
		db = txs[0]
	}
	return db.Save(&period).Error
}

func (r *studyTimetablePeriodRepository) UpdateStudyTimetablePeriod(period *models.StudyTimetablePeriod) error {
	return r.DB.Save(period).Error
}

func (r *studyTimetablePeriodRepository) DeleteStudyTimetablePeriod(period models.StudyTimetablePeriod) error {
	return r.DB.Unscoped().Delete(&period).Error
}

func (r *studyTimetablePeriodRepository) FindStudyTimetablePeriodByNumber(number string) (models.StudyTimetablePeriod, error) {
	var period models.StudyTimetablePeriod

	if err := r.DB.Where("number = ?", number).First(&period).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.StudyTimetablePeriod{}, nil
		}
		return models.StudyTimetablePeriod{}, err
	}

	return period, nil
}

func (r *studyTimetablePeriodRepository) FindStudyTimetablePeriods(studyTimetableNumber, streamNumber, day string) ([]models.StudyTimetablePeriod, error) {
	var periods []models.StudyTimetablePeriod

	query := r.DB
	if studyTimetableNumber != "" {
		query = query.Where("study_timetable_number = ?", studyTimetableNumber)
	}
	if streamNumber != "" {
		query = query.Where("stream_number = ?", streamNumber)
	}
	if day != "" {
		query = query.Where("day = ?", day)
	}

	if err := query.Find(&periods).Error; err != nil {
		return nil, err
	}

	return periods, nil
}

func (r *studyTimetablePeriodRepository) ExistsStudyTimetablePeriodByDayAndTime(streamNumber, studyPeriodNumber, day, startTime, endTime string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.StudyTimetablePeriod{}).
		Joins("JOIN study_timetables ON study_timetables.number = study_timetable_periods.study_timetable_number").
		Where("study_timetable_periods.stream_number = ?", streamNumber).
		Where("study_timetables.study_period_number = ?", studyPeriodNumber).
		Where("study_timetable_periods.day = ?", day).
		Where("study_timetable_periods.start_time = ?", startTime).
		Where("study_timetable_periods.end_time = ?", endTime).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
