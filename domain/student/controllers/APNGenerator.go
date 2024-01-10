package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
)

func CreateAPN() (string, error) {
	year := strconv.Itoa(time.Now().Year())
	month := time.Now().Format("01")

	saveAPN := models.APN{
		Year:  year,
		Month: month,
	}

	if err := config.DB.Create(&saveAPN).Error; err != nil {
		return "", err
	}

	savedID := saveAPN.ID + 1
	result := strings.Join([]string{saveAPN.Year, "01", fmt.Sprintf("%04s", strconv.FormatUint(uint64(savedID), 10))}, "")

	return result, nil
}
