package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
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
