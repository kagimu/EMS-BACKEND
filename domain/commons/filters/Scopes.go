package filters

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"github.com/gin-gonic/gin"
)

func GetUserScopedItems(ctx *gin.Context, category string) []string {
	var items = make([]string, 0)
	db := config.DB

	var st []models.StreamTeacher
	err := db.Where("user_number = ?", ctx.Request.Header["User-Number"][0]).First(&st).Error
	if err != nil {
		return items
	}

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

	// remove duplicates from the items slice using range
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i] == items[j] {
				items = append(items[:j], items[j+1:]...)
				j--
			}
		}
	}

	return items
}
