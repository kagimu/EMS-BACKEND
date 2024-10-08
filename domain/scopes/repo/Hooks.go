package repo

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"github.com/gin-gonic/gin"
)

func GetUserScopedItems(ctx *gin.Context, category string) []string {
	var items = make([]string, 0)
	db := config.DB
	user, _ := commons.GetUserFromHeader(ctx)
	var st []models.StreamTeacher
	err := db.Where("user_number = ?", user.Number).Find(&st).Error
	if err != nil {
		return items
	}

	for _, streamTeacher := range st {
		if category == "school" {
			items = append(items, streamTeacher.SchoolNumber)
		} else if category == "curriculum" {
			items = append(items, streamTeacher.CurriculumNumber)
		}

		if category == "campus" {
			items = append(items, streamTeacher.CampusNumber)
		}

		if category == "campusClass" {
			items = append(items, streamTeacher.CampusClassNumber)
		}

		if category == "stream" {
			items = append(items, streamTeacher.StreamNumber)
		}

		if category == "assessmentArea" {
			items = append(items, streamTeacher.AssessmentAreaNumber)
		}
	}

	return items
}
