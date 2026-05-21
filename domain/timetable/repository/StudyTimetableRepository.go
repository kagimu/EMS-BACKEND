package repository

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"gorm.io/gorm"
)

type StudyTimetableRepository interface {
	SaveStudyTimetable(timetable models.StudyTimetable, txs ...*gorm.DB) error
	UpdateStudyTimetable(timetable *models.StudyTimetable) error
	DeleteStudyTimetable(timetable models.StudyTimetable) error
	FindStudyTimetableByNumber(number string) (models.StudyTimetable, error)
	FindStudyTimetableByID(id uint) (models.StudyTimetable, error)
	FindAllStudyTimetables(studyPeriodNumber, userNumber, status, schoolNumber, levelNumber, campusNumber, campusCurriculumNumber, academicYearNumber, campusClassNumber, classNumber, streamNumber string) ([]models.StudyTimetable, error)
	ExistsStudyTimetable(schoolNumber, campusNumber, levelNumber, campusClassNumber, streamNumber, studyPeriodNumber string) (bool, error)
}

type studyTimetableRepository struct {
	DB *gorm.DB
}

func NewStudyTimetableRepository(db *gorm.DB) StudyTimetableRepository {
	return &studyTimetableRepository{DB: db}
}

func (r *studyTimetableRepository) SaveStudyTimetable(timetable models.StudyTimetable, txs ...*gorm.DB) error {
	db := r.DB
	if len(txs) > 0 {
		db = txs[0]
	}
	return db.Save(&timetable).Error
}

func (r *studyTimetableRepository) UpdateStudyTimetable(timetable *models.StudyTimetable) error {
	return r.DB.Save(timetable).Error
}

func (r *studyTimetableRepository) DeleteStudyTimetable(timetable models.StudyTimetable) error {
	return r.DB.Unscoped().Delete(&timetable).Error
}

func (r *studyTimetableRepository) FindStudyTimetableByNumber(number string) (models.StudyTimetable, error) {
	var timetable models.StudyTimetable

	if err := r.DB.Where("number = ?", number).First(&timetable).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.StudyTimetable{}, nil
		}
		return models.StudyTimetable{}, err
	}

	return timetable, nil
}

func (r *studyTimetableRepository) FindStudyTimetableByID(id uint) (models.StudyTimetable, error) {
	var timetable models.StudyTimetable

	if err := r.DB.Where("id = ?", id).First(&timetable).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.StudyTimetable{}, nil
		}
		return models.StudyTimetable{}, err
	}

	return timetable, nil
}

func (r *studyTimetableRepository) FindAllStudyTimetables(studyPeriodNumber, userNumber, status, schoolNumber, levelNumber, campusNumber, campusCurriculumNumber, academicYearNumber, campusClassNumber, classNumber, streamNumber string) ([]models.StudyTimetable, error) {
	var timetables []models.StudyTimetable

	query := r.DB
	if studyPeriodNumber != "" {
		query = query.Where("study_period_number = ?", studyPeriodNumber)
	}
	if userNumber != "" {
		query = query.Where("user_number = ?", userNumber)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if schoolNumber != "" {
		query = query.Where("school_number = ?", schoolNumber)
	}
	if levelNumber != "" {
		query = query.Where("level_number = ?", levelNumber)
	}
	if campusNumber != "" {
		query = query.Where("campus_number = ?", campusNumber)
	}
	if campusCurriculumNumber != "" {
		query = query.Where("campus_curriculum_number = ?", campusCurriculumNumber)
	}
	if academicYearNumber != "" {
		query = query.Where("academic_year_number = ?", academicYearNumber)
	}
	if campusClassNumber != "" {
		query = query.Where("campus_class_number = ?", campusClassNumber)
	}
	if classNumber != "" {
		// class_number is not a direct column on study_timetables.
		// Resolve: campus_classes.curriculum_class_number → curriculum_classes.class_number
		subQuery := r.DB.Table("campus_classes").
			Select("campus_classes.number").
			Joins("JOIN curriculum_classes ON curriculum_classes.number = campus_classes.curriculum_class_number").
			Where("curriculum_classes.class_number = ?", classNumber)
		query = query.Where("campus_class_number IN (?)", subQuery)
	}
	if streamNumber != "" {
		query = query.Where("stream_number = ?", streamNumber)
	}

	err := query.Find(&timetables).Error
	if err != nil {
		return nil, err
	}

	return timetables, nil
}

func (r *studyTimetableRepository) ExistsStudyTimetable(schoolNumber, campusNumber, levelNumber, campusClassNumber, streamNumber, studyPeriodNumber string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.StudyTimetable{}).
		Where("school_number = ?", schoolNumber).
		Where("campus_number = ?", campusNumber).
		Where("level_number = ?", levelNumber).
		Where("campus_class_number = ?", campusClassNumber).
		Where("stream_number = ?", streamNumber).
		Where("study_period_number = ?", studyPeriodNumber).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
