package repo

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"github.com/gin-gonic/gin"
)

func GetUserScopedItems(ctx *gin.Context, category string) []string {
	var items = make([]string, 0)
	db := config.DB
	user, _ := commons.GetUserFromHeader(ctx)
	fmt.Println("User: ", user.FirstName)

	var st []models.StreamTeacher
	err := db.Where("user_number = ?", user.Number).First(&st).Error

	if err != nil {
		fmt.Println(fmt.Sprintf("Error: %s", err.Error()))
		return items
	}

	fmt.Println(fmt.Sprintf("length is %d", len(st)))
	for _, streamTeacher := range st {
		switch category {
		case "school":
			items = append(items, streamTeacher.SchoolNumber)
		case "curriculum":
			items = append(items, streamTeacher.CurriculumNumber)
		case "campus":
			items = append(items, streamTeacher.CampusNumber)
		case "campusClass":
			items = append(items, streamTeacher.CampusClassNumber)
		case "stream":
			items = append(items, streamTeacher.StreamNumber)
		case "assessmentArea":
			items = append(items, streamTeacher.AssessmentAreaNumber)
		}
	}

	return items
}
