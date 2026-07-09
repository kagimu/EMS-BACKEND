package repository

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"gorm.io/gorm"
)

type StaffScopeRepository interface {
	ReplaceStaffScopes(staffNumber, schoolNumber string, scopes []models.StaffScope) error
	GetStaffScopes(staffNumber, schoolNumber string) ([]models.StaffScope, error)
	GetStaffScopesByUserNumber(userNumber string) ([]models.StaffScope, error)
	ClearStaffScope(staffNumber, schoolNumber string) error
}

type staffScopeRepo struct {
	DB *gorm.DB
}

func NewStaffScopeRepository(db *gorm.DB) StaffScopeRepository {
	return &staffScopeRepo{DB: db}
}

func (s *staffScopeRepo) ClearStaffScope(staffNumber, schoolNumber string) error {
	query := s.DB.Where("staff_number = ?", staffNumber)
	if schoolNumber != "" {
		query = query.Where("school_number = ?", schoolNumber)
	}
	return query.Delete(&models.StaffScope{}).Error
}

func (s *staffScopeRepo) ReplaceStaffScopes(staffNumber, schoolNumber string, scopes []models.StaffScope) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("staff_number = ? AND school_number = ?", staffNumber, schoolNumber).Delete(&models.StaffScope{}).Error; err != nil {
			return err
		}
		if len(scopes) > 0 {
			if err := tx.Create(&scopes).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *staffScopeRepo) GetStaffScopes(staffNumber, schoolNumber string) ([]models.StaffScope, error) {
	var scopes []models.StaffScope
	query := s.DB.Where("staff_number = ?", staffNumber)
	if schoolNumber != "" {
		query = query.Where("school_number = ?", schoolNumber)
	}
	err := query.Find(&scopes).Error
	return scopes, err
}

func (s *staffScopeRepo) GetStaffScopesByUserNumber(userNumber string) ([]models.StaffScope, error) {
	var scopes []models.StaffScope
	err := s.DB.Where("user_number = ? OR staff_number = ?", userNumber, userNumber).Find(&scopes).Error
	return scopes, err
}

