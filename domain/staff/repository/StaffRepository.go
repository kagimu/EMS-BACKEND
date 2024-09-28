package repository

import (
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"gorm.io/gorm"
)

type StaffRepository interface {
	AttachStaffToStream(streamTeacher models.StreamTeacher) error
	DetachStaffFromStream(staffNumber string, streamNumber string) error
	GetStaffOnStream(streamNumber string) ([]userModel.User, error)
}

type StaffRepo struct {
	DB *gorm.DB
}

func (s StaffRepo) AttachStaffToStream(streamTeacher models.StreamTeacher) error {
	return s.DB.Create(&streamTeacher).Error
}

func (s StaffRepo) DetachStaffFromStream(staffNumber string, streamNumber string) error {
	return s.DB.Where("user_number = ? AND stream_number = ?", staffNumber, streamNumber).Delete(&models.StreamTeacher{}).Error
}

func (s StaffRepo) GetStaffOnStream(streamNumber string) ([]userModel.User, error) {
	var staff []models.StreamTeacher
	var users = make([]userModel.User, 0)
	err := s.DB.Where("stream_number = ?", streamNumber).Find(&staff).Error
	if err != nil {
		return users, err
	}
	for _, v := range staff {
		var user userModel.User
		err := s.DB.Where("number = ?", v.UserNumber).First(&user).Error
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, nil

}

func NewStaffRepository(db *gorm.DB) StaffRepository {
	return &StaffRepo{DB: db}
}
