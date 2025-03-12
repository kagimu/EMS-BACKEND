package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	itemModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	"gorm.io/gorm"
)

func GetItemByNumber(itemNumber string, tx ...*gorm.DB) (itemModels.Item, error) {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	var item itemModels.Item

	if err := db.Where("number = ?", itemNumber).First(&item).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return item, nil
		}

		return item, err
	}

	return item, nil
}

func FetchAllServiceItems(schoolNumber string, tx ...*gorm.DB) ([]itemModels.Item, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var items []itemModels.Item
	if err := db.
		Joins("JOIN services sv ON items.number = sv.item_number").
		Scopes(filters.ByAnyField("items.school_number", schoolNumber)).Find(&items).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []itemModels.Item{}, nil
		}
		return []itemModels.Item{}, err
	}

	return items, nil
}

func FetchAllCampusItems(itemNumber string, tx ...*gorm.DB) ([]campusModel.Campus, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var campuses []campusModel.Campus
	if err := db.
		Joins("JOIN campus_items ON campuses.number = campus_items.campus_number").
		Scopes(filters.ByAnyField("campus_items.item_number", itemNumber)).Find(&campuses).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []campusModel.Campus{}, nil
		}
		return []campusModel.Campus{}, err
	}

	return campuses, nil
}
