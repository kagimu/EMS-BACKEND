package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type StudentRepository interface {
	FindAll() ([]studentModel.Student, error)
	GetStudentByApplicationNumber(studentNumber string) (studentModel.Student, error)
	Delete(studentNumber string, tx *gorm.DB) error
}

type StudentRepositoryImpl struct {
	DB *gorm.DB
}

func (s StudentRepositoryImpl) GetStudentByApplicationNumber(applicationNumber string) (studentModel.Student, error) {
	var student studentModel.Student

	if err := s.DB.Where("application_number = ?", applicationNumber).First(&student).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return student, nil
		}

		return student, err
	}

	return student, nil
}

func (s StudentRepositoryImpl) FindAll() ([]studentModel.Student, error) {
	//TODO implement me
	panic("implement me")
}

func (s StudentRepositoryImpl) Delete(studentNumber string, txn *gorm.DB) error {
	return s.DB.Where("number = ?", studentNumber).Delete(&studentModel.Student{}).Error
}

func NewStudentRepositoryImpl(db *gorm.DB) StudentRepository {
	return &StudentRepositoryImpl{DB: db}
}

func FindStudentByUserNumber(userNumber string, txn ...*gorm.DB) (studentModel.Student, error) {

	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var student studentModel.Student

	if err := db.Where("user_number = ?", userNumber).First(&student).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return studentModel.Student{}, nil
		}

		return studentModel.Student{}, err
	}

	return student, nil

}

func FindStudentByNumber(studentNumber string) (studentModel.Student, error) {

	db := config.DB
	var student studentModel.Student

	if err := db.Where("number = ?", studentNumber).First(&student).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return student, nil
		}

		return student, err
	}

	return student, nil

}
