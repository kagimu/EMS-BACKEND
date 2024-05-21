package repository

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/statistics/api/responses"
	"github.com/gin-gonic/gin"
)

func FetchBasicStatistics(ctx *gin.Context) (responses.BasicStatistics, error) {
	var usersCount int
	config.DB.Raw("SELECT COUNT(*) FROM users").Scan(&usersCount)

	var classesCount int
	config.DB.Raw("SELECT COUNT(*) FROM campus_classes").Scan(&classesCount)

	var parentsCount int
	config.DB.Raw("SELECT COUNT(*) FROM parents").Scan(&parentsCount)

	var enrolledStudentsCount int
	config.DB.Raw("SELECT COUNT(*) FROM enrollments WHERE status = 'Active'").Scan(&enrolledStudentsCount)

	var teachersCount int
	config.DB.Raw("SELECT COUNT(*) FROM users WHERE type = 'Staff'").Scan(&teachersCount)

	return responses.BasicStatistics{
		TotalTeachers: teachersCount,
		TotalStudents: enrolledStudentsCount,
		TotalClasses:  classesCount,
		TotalCampuses: 0,
		TotalParents:  parentsCount,
	}, nil
}
