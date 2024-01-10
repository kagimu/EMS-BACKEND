package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
)

func CreateADN() (string, error) {
	year := strconv.Itoa(time.Now().Year())
	month := time.Now().Format("01")

	saveADN := models.ADN{
		Year:  year,
		Month: month,
	}

	if err := config.DB.Create(&saveADN).Error; err != nil {
		return "", err
	}

	savedID := saveADN.ID + 1
	result := strings.Join([]string{saveADN.Year, "02", fmt.Sprintf("%04s", strconv.FormatUint(uint64(savedID), 10))}, "")

	return result, nil
}
