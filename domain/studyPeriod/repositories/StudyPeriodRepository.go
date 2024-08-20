package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"gorm.io/gorm"
)

func GetStudyPeriodByNumber(studyPeriodNumber string) (studyPeriodModel.StudyPeriod, error) {
	db := config.DB
	var studyPeriod studyPeriodModel.StudyPeriod
	if err := db.Where("number = ?", studyPeriodNumber).First(&studyPeriod).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return studyPeriod, nil
		}

		return studyPeriod, err
	}
	return studyPeriod, nil
}

func GetStudyPeriodByCode(academicYearNumber string, code int64) (studyPeriodModel.StudyPeriod, error) {
	db := config.DB
	var studyPeriod studyPeriodModel.StudyPeriod
	if err := db.Where("academic_year_number = ? AND code = ?", academicYearNumber, code).First(&studyPeriod).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return studyPeriod, nil
		}

		return studyPeriod, err
	}
	return studyPeriod, nil
}

func SaveStudyPeriod(studyPeriod studyPeriodModel.StudyPeriod, trn ...*gorm.DB) (studyPeriodModel.StudyPeriod, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Save(&studyPeriod).Error; err != nil {
		return studyPeriodModel.StudyPeriod{}, err
	}

	return studyPeriod, nil
}

func RemoveStudyPeriod(studyPeriod studyPeriodModel.StudyPeriod, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Delete(&studyPeriod).Error
}

func FindStudyPeriod(academicYearNumber, calendarEventNumber, studyPeriodNameNumber, classNumber string, trn ...*gorm.DB) ([]studyPeriodModel.StudyPeriod, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var studyPeriods []studyPeriodModel.StudyPeriod

	if err := db.
		Table("study_periods").
		Select("DISTINCT study_periods.*").
		Joins("join class_study_periods on class_study_periods.study_period_number = study_periods.number").
		Scopes(filters.ByAcademicYearNumber(academicYearNumber), filters.ByAnyField("class_study_periods.class_number", classNumber), filters.ByAnyField("study_period_name_number", studyPeriodNameNumber), filters.ByAnyField("calendar_event_number", calendarEventNumber)).
		Order("code asc").
		Find(&studyPeriods).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return studyPeriods, nil
		}

		return studyPeriods, err
	}

	return studyPeriods, nil
}
