package repository

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studentDetails"
	"gorm.io/gorm"
)

type StudentPhotoRepo interface {
	CreateStudentPhoto(studentPhoto *studentDetails.StudentPhoto) error
	GetStudentPhoto(studentNumber string) (*studentDetails.StudentPhoto, error)
	GetStudentPhotosByStudentNumber(studentNumber string) ([]studentDetails.StudentPhoto, error)
}

type StudentPhotoRepoImpl struct {
	db *gorm.DB
}

func NewStudentPhotoRepo(db *gorm.DB) *StudentPhotoRepoImpl {
	return &StudentPhotoRepoImpl{db: db}
}

func (s *StudentPhotoRepoImpl) CreateStudentPhoto(studentPhoto *studentDetails.StudentPhoto) error {
	err := s.db.Create(studentPhoto).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *StudentPhotoRepoImpl) GetStudentPhoto(studentNumber string) (*studentDetails.StudentPhoto, error) {
	var studentPhoto studentDetails.StudentPhoto
	err := s.db.Where("student_number = ?", studentNumber).First(&studentPhoto).Error
	if err != nil {
		return nil, err
	}
	return &studentPhoto, nil
}

func (s *StudentPhotoRepoImpl) GetStudentPhotosByStudentNumber(studentNumber string) ([]studentDetails.StudentPhoto, error) {
	var studentPhotos []studentDetails.StudentPhoto
	err := s.db.Where("student_number = ?", studentNumber).Find(&studentPhotos).Error
	if err != nil {
		return nil, err
	}
	return studentPhotos, nil
}
