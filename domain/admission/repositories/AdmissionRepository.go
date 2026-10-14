package repositories

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	admissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

type AdmissionRepository interface {
	FindAdmissionByStudentNumber(studentNumber string, txn ...*gorm.DB) (admissionModel.Admission, error)
	SaveAdmission(admission admissionModel.Admission, txn ...*gorm.DB) error
	AllAdmissions(campusNumber, studyPeriodNumber, streamNumber, schoolNumber, campusClassNumber, search, academicYearNumber, classNumber, levelNumber, curriculumNumber string, user *userModel.User, paginateFunc func(db *gorm.DB) *gorm.DB, isMine ...bool) ([]admissionModel.Admission, error)
	FindAdmissionsByStreamNumber(streamNumber string, txn ...*gorm.DB) ([]admissionModel.Admission, error)
}

type admissionRepository struct {
}

func (admissionRepository *admissionRepository) AllAdmissions(campusNumber, studyPeriodNumber, streamNumber, schoolNumber, campusClassNumber, search, academicYearNumber, classNumber, levelNumber, curriculumNumber string, user *userModel.User, paginateFunc func(db *gorm.DB) *gorm.DB, isMine ...bool) ([]admissionModel.Admission, error) {

	db := config.DB
	mine := len(isMine) > 0 && isMine[0]

	var admissions []admissionModel.Admission

	if err := db.Model(&admissionModel.Admission{}).
		Joins("JOIN students ON students.number = admissions.student_number").
		Joins("JOIN applications ON applications.number = students.application_number").
		Scopes(
			filters.ByAnyField("admissions.study_period_number", studyPeriodNumber),
			filters.ByAnyField("admissions.campus_number", campusNumber),
			filters.ByAnyField("admissions.stream_number", streamNumber),
			filters.ByAnyField("admissions.school_number", schoolNumber),
			filters.ByAnyField("admissions.campus_class_number", campusClassNumber),
			filters.ByAnyField("admissions.academic_year_number", academicYearNumber),
			filters.BySearchingAnyField([]string{"students.adn", "students.apn", "students.surname", "students.second_name", "students.other_name"}, search),
			func(db *gorm.DB) *gorm.DB {
				q := db
				if classNumber != "" {
					q = q.Where("admissions.campus_class_number IN (SELECT number FROM campus_classes WHERE deleted_at IS NULL AND curriculum_class_number IN (SELECT number FROM curriculum_classes WHERE deleted_at IS NULL AND class_number = ?))", classNumber)
				}
				if levelNumber != "" {
					q = q.Where("admissions.campus_class_number IN (SELECT number FROM campus_classes WHERE deleted_at IS NULL AND campus_level_number = ?)", levelNumber)
				}
				if curriculumNumber != "" {
					q = q.Where("admissions.campus_class_number IN (SELECT number FROM campus_classes WHERE deleted_at IS NULL AND campus_level_number IN (SELECT number FROM campus_levels WHERE deleted_at IS NULL AND campus_curriculum_number = ?))", curriculumNumber)
				}
				return q
			},
			func(db *gorm.DB) *gorm.DB {
				return db.Where("admissions.status IN (?)", []string{constants.StatusPending, constants.StatusAdmitted})
			},
			filters.ByParent("applications.user_number", user, mine),
			paginateFunc,
		).Find(&admissions).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return admissions, nil
		}

		return admissions, err
	}

	return admissions, nil
}

func (admissionRepository *admissionRepository) SaveAdmission(admission admissionModel.Admission, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	return db.Save(&admission).Error
}

func (admissionRepository *admissionRepository) FindAdmissionByStudentNumber(studentNumber string, txn ...*gorm.DB) (admissionModel.Admission, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var admission admissionModel.Admission

	if err := db.Where("student_number = ?", studentNumber).First(&admission).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return admission, nil
		}

		return admission, err
	}

	return admission, nil
}

func (admissionRepository *admissionRepository) FindAdmissionsByStreamNumber(streamNumber string, txn ...*gorm.DB) ([]admissionModel.Admission, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var admissions []admissionModel.Admission

	if err := db.Where("stream_number = ?", streamNumber).Find(&admissions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return admissions, nil
		}
		return admissions, err
	}

	return admissions, nil
}

func NewAdmissionRepository() AdmissionRepository {
	return &admissionRepository{}
}
