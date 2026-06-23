package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	departmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/departments/models"
	"gorm.io/gorm"
)

func FindDepartmentByNumber(departmentNumber string, trn ...*gorm.DB) (departmentModel.Department, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var department departmentModel.Department

	if err := db.Where("number = ?", departmentNumber).First(&department).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return department, nil
		}
		return department, err
	}

	return department, nil
}
