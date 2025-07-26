package repository

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"gorm.io/gorm"
)

type EnrolmentRepository interface {
	GetStudentEnrolment(studentNumber, studyPeriodNumber, campusClassNumber, status string, txn ...*gorm.DB) (models.Enrollment, error)
	AllEnrolments(studentNumber, studyPeriodNumber, campusClassNumber, campusNumber, status string, txn ...*gorm.DB) ([]models.Enrollment, error)
	SaveEnrolment(enrollment models.Enrollment, txn ...*gorm.DB) error
	PaginatedEnrolments(schoolNumber, streamNumber, campusNumber, campusClassNumber, studyPeriodNumber, academicYearNumber, status string, user *userModel.User, page func(db *gorm.DB) *gorm.DB) ([]models.Enrollment, error)
}

type enrolmentRepository struct {
}

func (enrolmentRepository *enrolmentRepository) PaginatedEnrolments(schoolNumber, streamNumber, campusNumber, campusClassNumber, studyPeriodNumber, academicYearNumber, status string, user *userModel.User, page func(db *gorm.DB) *gorm.DB) ([]models.Enrollment, error) {
	var enrolledStudents []models.Enrollment
	if err := config.DB.
		Joins("JOIN study_periods ON study_periods.number = enrollments.study_period_number").
		Joins("JOIN students ON students.number = enrollments.student_number").
		Joins("JOIN applications ON applications.number = students.application_number").
		Scopes(
			filters.ByAnyField("enrollments.school_number", schoolNumber),
			filters.ByAnyField("enrollments.stream_number", streamNumber),
			filters.ByAnyField("enrollments.campus_number", campusNumber),
			filters.ByAnyField("enrollments.campus_class_number", campusClassNumber),
			filters.ByAnyField("enrollments.study_period_number", studyPeriodNumber),
			filters.ByParent("applications.user_number", user),
			filters.ByAnyField("enrollments.status", status),
			filters.ByAnyField("study_periods.academic_year_number", academicYearNumber),
			page,
		).
		Preload("School").Preload("CampusClass").Preload("Campus").Preload("Student").Preload("StudyPeriod").Find(&enrolledStudents).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Enrollment{}, nil
		}

		return nil, err
	}

	return enrolledStudents, nil
}

func (enrolmentRepository *enrolmentRepository) AllEnrolments(studentNumber, studyPeriodNumber, campusClassNumber, campusNumber, status string, txn ...*gorm.DB) ([]models.Enrollment, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var enrolment []models.Enrollment
	if err := db.
		Scopes(filters.ByStatus(status), filters.ByCampusNumber(campusNumber), filters.ByStudyPeriodNumber(studyPeriodNumber), filters.ByStudentNumber(studentNumber), filters.ByCampusClassNumber(campusClassNumber)).
		Find(&enrolment).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Enrollment{}, nil
		}

		return []models.Enrollment{}, err
	}

	return enrolment, nil
}

func (enrolmentRepository *enrolmentRepository) SaveEnrolment(enrollment models.Enrollment, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	return db.Save(&enrollment).Error
}

func (enrolmentRepository *enrolmentRepository) GetStudentEnrolment(studentNumber, studyPeriodNumber, campusClassNumber, status string, txn ...*gorm.DB) (models.Enrollment, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var enrolment models.Enrollment
	if err := db.
		Scopes(filters.ByStatus(status), filters.ByCampusClassNumber(campusClassNumber)).
		Where("student_number = ? AND study_period_number = ?", studentNumber, studyPeriodNumber).
		First(&enrolment).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Enrollment{}, nil
		}

		return models.Enrollment{}, err
	}

	return enrolment, nil
}

func NewEnrolmentRepository() EnrolmentRepository {
	return &enrolmentRepository{}
}
