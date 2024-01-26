package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
)

func CreateADN(apnStatus string) (string, error) {
	year := strconv.Itoa(time.Now().Year())
	month := time.Now().Format("01")

	if apnStatus == "Archive" {
		return "Status is already archived", nil
	}

	saveADN := models.ADN{
		Year:  year,
		Month: month,
	}

	if err := config.DB.Create(&saveADN).Error; err != nil {
		return "", err
	}

	savedID := saveADN.ID
	result := fmt.Sprintf("%s02%04d", saveADN.Year, savedID)

	return result, nil
}
