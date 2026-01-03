package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/controllers/api"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

var sourceOfSchoolService = services.NewSourceOfSchoolService()

func NewSourceOfSchoolInfo(c *gin.Context) {

	var sourceOfSchoolInfo api.SourceOfSchoolRequest

	err := c.ShouldBindJSON(&sourceOfSchoolInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = sourceOfSchoolInfo.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userNumber := c.Request.Header["User-Number"][0]

	entityNumber, err := commons.GetUserEntityNumber(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	data, err := sourceOfSchoolService.NewSourceOfSchool(sourceOfSchoolInfo, userNumber, entityNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "source of school information created successfully",
		"data":    data,
	})

}

func GetSourcesOfSchoolInfo(c *gin.Context) {

	entityNumber := c.Query("entity_number")

	if entityNumber == "" {
		entityNumber, _ = commons.GetUserEntityNumber(c)
	}

	data, err := sourceOfSchoolService.FindAllSources(entityNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": data})
}

func EditSourceOfSchoolInfo(c *gin.Context) {

	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var sourceOfSchoolInfo api.SourceOfSchoolRequest
	if err := c.ShouldBindJSON(&sourceOfSchoolInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input- " + err.Error()})
		return
	}

	err = sourceOfSchoolInfo.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	sourceOfSchoolInfoNumber := c.Param("source_of_school_info_number")

	sourceOfSchoolInfoModel, err := sourceOfSchoolService.Update(sourceOfSchoolInfo, sourceOfSchoolInfoNumber, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "source of school info edited successfully", "data": sourceOfSchoolInfoModel})
}

func DeleteSourceOfSchoolInfo(c *gin.Context) {

	entityNumber, err := commons.GetUserEntityNumber(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	sourceOfSchoolInfoNumber := c.Param("source_of_school_info_number")

	err = sourceOfSchoolService.DeleteSource(entityNumber, sourceOfSchoolInfoNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Source of school info deleted successfully"})
}
