package filters

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"github.com/gin-gonic/gin"
)

func GetUserScopedItems(ctx *gin.Context, category string, userNumber ...string) []string {
	var items = make([]string, 0)
	db := config.DB
	var uNumber string
	if len(userNumber) > 0 {
		uNumber = userNumber[0]
	} else if len(ctx.Request.Header["User-Number"]) > 0 {
		uNumber = ctx.Request.Header["User-Number"][0]
	}

	if uNumber == "" {
		return items
	}

	var st []models.StreamTeacher
	db.Where("user_number = ?", uNumber).Find(&st)

	for _, streamTeacher := range st {
		switch category {
		case "school":
			items = append(items, streamTeacher.SchoolNumber)
		case "curriculum":
			items = append(items, streamTeacher.CurriculumNumber)
		case "campus":
			items = append(items, streamTeacher.CampusNumber)
		case "campusClass", "class":
			items = append(items, streamTeacher.CampusClassNumber)
		case "stream":
			items = append(items, streamTeacher.StreamNumber)
		case "assessmentArea":
			items = append(items, streamTeacher.AssessmentAreaNumber)
		}
	}

	var staffList []models.Staff
	db.Where("user_number = ? OR number = ?", uNumber, uNumber).Find(&staffList)
	numbers := []string{uNumber}
	for _, s := range staffList {
		if s.Number != "" {
			numbers = append(numbers, s.Number)
		}
		if s.UserNumber != "" {
			numbers = append(numbers, s.UserNumber)
		}
	}

	var staffScopes []models.StaffScope
	db.Where("user_number IN (?) OR staff_number IN (?)", numbers, numbers).Find(&staffScopes)

	for _, scope := range staffScopes {
		switch category {
		case "school":
			items = append(items, scope.SchoolNumber)
		case "curriculum":
			items = append(items, scope.CurriculumNumber)
		case "campus":
			items = append(items, scope.CampusNumber)
		case "campusLevel":
			items = append(items, scope.CampusLevelNumber)
		case "campusClass", "class":
			items = append(items, scope.ClassNumber)
		case "stream":
			items = append(items, scope.StreamNumber)
		case "assessmentArea":
			items = append(items, scope.AssessmentAreaNumber)
		}
	}

	cleanItems := make([]string, 0, len(items))
	seen := make(map[string]bool)
	for _, val := range items {
		if val != "" && !seen[val] {
			seen[val] = true
			cleanItems = append(cleanItems, val)
		}
	}

	return cleanItems
}
