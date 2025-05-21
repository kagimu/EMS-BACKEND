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

func FindStudyPeriod(academicYearNumber, calendarEventNumber, studyPeriodNameNumber, classNumber, schoolNumber string, trn ...*gorm.DB) ([]studyPeriodModel.StudyPeriod, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var studyPeriods []studyPeriodModel.StudyPeriod

	if err := db.
		Joins("JOIN study_period_names ON study_period_names.number = study_periods.study_period_name_number").
		Joins("JOIN class_study_periods ON class_study_periods.study_period_number = study_periods.number").
		Scopes(filters.ByAcademicYearNumber(academicYearNumber), filters.ByAnyField("", schoolNumber), filters.ByAnyField("class_study_periods.class_number", classNumber), filters.ByAnyField("study_period_name_number", studyPeriodNameNumber), filters.ByAnyField("calendar_event_number", calendarEventNumber)).
		Order("code asc").
		Distinct("study_periods.number", "study_periods.code", "study_periods.study_period_name_number", "study_periods.study_period_start_date", "study_periods.study_period_end_date", "study_periods.academic_year_number", "study_periods.calendar_event_number", "study_periods.user_number", "study_periods.next_term_begins_on").
		Find(&studyPeriods).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return studyPeriods, nil
		}

		return studyPeriods, err
	}

	return studyPeriods, nil
}

func GetDuplicateClassStudyPeriod(classNumber, academicYearNumber, studyPeriodNameNumber string, trn ...*gorm.DB) (studyPeriodModel.StudyPeriod, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var studyPeriod studyPeriodModel.StudyPeriod

	if err := db.
		Joins("JOIN class_study_periods ON class_study_periods.study_period_number = study_periods.number").
		Where("academic_year_number = ? AND class_number = ? AND study_period_name_number = ?", academicYearNumber, classNumber, studyPeriodNameNumber).
		Find(&studyPeriod).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return studyPeriod, nil
		}

		return studyPeriod, err
	}

	return studyPeriod, nil
}
