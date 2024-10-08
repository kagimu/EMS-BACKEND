package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateStream(stream *models.Stream, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}
	if err := db.Create(&stream).Error; err != nil {
		return err
	}
	return nil
}

func FetchStreamByCampusClassNumber(campusClassNumber string, ctx *gin.Context, trn ...*gorm.DB) ([]models.Stream, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var streams []models.Stream

	if err := db.Scopes(filters.ByScopedItems("stream", ctx)).Where("campus_class_number = ?", campusClassNumber).Find(&streams).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return streams, nil
		}

		return streams, err

	}

	return streams, nil
}

func GetStreamByNumber(streamNumber string, trn ...*gorm.DB) (models.Stream, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var stream models.Stream

	if err := db.Where("number = ?", streamNumber).Find(&stream).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return stream, nil
		}

		return stream, err

	}

	return stream, nil
}
