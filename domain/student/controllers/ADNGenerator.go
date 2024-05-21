package controllers

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
)

func CreateADN(apnStatus string, tx *gorm.DB) (string, error) {
	year := strconv.Itoa(time.Now().Year())
	month := time.Now().Format("02") // Get zero-padded month

	if apnStatus == "Archive" {
		return "", nil
	} else {
		saveADN := models.ADN{
			Year:  year,
			Month: month,
		}

		if err := tx.Create(&saveADN).Error; err != nil {
			return "", err
		}

		savedID := saveADN.ID
		result := fmt.Sprintf("%s02%04d", saveADN.Year, savedID)

		return result, nil
	}

}
