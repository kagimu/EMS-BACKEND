package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMenuItemByNumber(menuItemNumber string) (models.MenuItem, error) {

	db := config.DB

	var menuItem models.MenuItem

	if err := db.Where("number = ?", menuItemNumber).First(&menuItem).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return menuItem, errors.New("invalid menu item number")
		}

		return menuItem, err
	}

	return menuItem, nil
}

func GetMenuItemByNameAndSchool(schoolNumber string, menuItemName string) (models.MenuItem, error) {

	db := config.DB

	var menuItem models.MenuItem

	if err := db.
		Where("school_number = ?", schoolNumber).
		Where("name = ?", menuItemName).
		First(&menuItem).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return menuItem, nil
		}

		return menuItem, err
	}

	return menuItem, nil
}

func GetMenuItemByCodeAndSchool(schoolNumber string, itemCode int64) (models.MenuItem, error) {
	db := config.DB

	var menuItem models.MenuItem

	if err := db.
		Where("school_number = ?", schoolNumber).
		Where("code = ?", itemCode).
		First(&menuItem).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return menuItem, nil
		}

		return menuItem, err
	}

	return menuItem, nil
}

func SaveMenuItem(menuItem models.MenuItem) error {

	db := config.DB

	if err := db.Save(&menuItem).Error; err != nil {
		return err
	}

	return nil

}

func GetPaginatedMenuItems(schoolNumber string, status string, paginatedResponse pagination.PaginatedResponse, c *gin.Context) ([]models.MenuItem, error) {

	var menuItems []models.MenuItem

	if err := config.DB.
		Scopes(filters.BySchoolNumber(schoolNumber), filters.ByStatus(status), pagination.
			Paginate(c, &paginatedResponse)).Preload("School").Preload("UnitOfMeasure").
		Find(&menuItems).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return menuItems, nil
		}

		return menuItems, err
	}

	return menuItems, nil
}

func GetAllMenuItems(schoolNumber string, status string) ([]models.MenuItem, error) {

	var menuItems []models.MenuItem

	if err := config.DB.
		Scopes(filters.BySchoolNumber(schoolNumber), filters.ByStatus(status)).Preload("School").Preload("UnitOfMeasure").
		Find(&menuItems).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return menuItems, nil
		}

		return menuItems, err
	}

	return menuItems, nil
}

func DeleteMenuItem(item models.MenuItem) error {

	db := config.DB

	if err := db.Delete(&item).Error; err != nil {
		return err
	}

	return nil
}
