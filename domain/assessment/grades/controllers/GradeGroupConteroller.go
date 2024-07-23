package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	uuid "github.com/satori/go.uuid"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateGradeGroup(ctx *gin.Context) {
	var err error
	// Implement the logic to create a new grade group
	var req requests.CreateGradeGroupRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input " + err.Error()})
		return
	}

	err = req.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var gradeGroup models.GradeGroup

	gradeGroup.Name = req.Name
	gradeGroup.Code = req.Code
	gradeGroup.Number = uuid.NewV4().String()
	gradeGroup.UserNumber = ctx.Request.Header["User-Number"][0]
	gradeGroup.SchoolNumber = req.SchoolNumber

	// Call the service to create a new grade group
	err = services.CreateGradeGroup(&gradeGroup)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
		return

	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Grade group created successfully"})

}

func FetchGradeGroups(ctx *gin.Context) {

	schoolNumber := ctx.Query("school_number")

	if schoolNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "School number is required"})
		return
	}

	// Implement the logic to fetch all grade groups
	gradeGroups, err := services.FetchGradeGroups(schoolNumber)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch grade groups"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"gradeGroups": gradeGroups})
}

func UpdateGradeGroup(ctx *gin.Context) {

	var req requests.UpdateGradeGroupRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	err := req.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	gradeGroup, err := services.FetchGradeGroupByNumber(req.Number)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Grade group not found"})
		return

	}

	gradeGroup.Name = req.Name
	gradeGroup.Code = req.Code
	gradeGroup.UserNumber = ctx.Request.Header["User-Number"][0]

	// Call the service to update the grade group
	err = services.UpdateGradeGroup(&gradeGroup)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update grade group"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Grade group updated successfully"})
}

func DeleteGradeGroup(ctx *gin.Context) {
	// Implement the logic to delete a grade group
	gradeGroupID := ctx.Param("number")

	err := services.DeleteGradeGroup(gradeGroupID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete grade group"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Grade group deleted successfully"})
}
