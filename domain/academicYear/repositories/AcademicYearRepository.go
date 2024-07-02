package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	academicYearModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	"gorm.io/gorm"
)

func GetAcademicYearByNumber(academicYearNumber string) (academicYearModels.AcademicYear, error) {
	db := config.DB
	var academicYear academicYearModels.AcademicYear
	if err := db.Where("number = ?", academicYearNumber).First(&academicYear).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return academicYear, nil
		}
		return academicYear, err
	}
	return academicYear, nil
}
