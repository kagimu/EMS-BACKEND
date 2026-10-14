package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/addons/models"
)

func FindAllModules() ([]models.Module, error) {
	var modules []models.Module
	if err := config.DB.Order("code asc").Find(&modules).Error; err != nil {
		return nil, err
	}
	return modules, nil
}

func FindModuleByCode(code string) (*models.Module, error) {
	var module models.Module
	if err := config.DB.Where("code = ?", code).First(&module).Error; err != nil {
		return nil, err
	}
	return &module, nil
}

func CreateModule(module *models.Module) error {
	return config.DB.Create(module).Error
}

func SaveModule(module *models.Module) error {
	return config.DB.Save(module).Error
}
