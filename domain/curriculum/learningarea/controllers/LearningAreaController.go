package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/api/responses"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

func NewLearningArea(context *gin.Context) {

	type LearningAreaInfo struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}

	var learningAreaInfo LearningAreaInfo
	err := context.BindJSON(&learningAreaInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	name := learningAreaInfo.Name
	code := learningAreaInfo.Code

	userNumber := context.Request.Header["User-Number"][0]

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code is required!"})
		return
	}

	codeNo, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var userData models.User
	//check if user is a valid user
	validateUser := config.DB.
		Where("number = ?", userNumber).
		First(&userData)
	if validateUser.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var duplicateLearningArea models2.LearningArea
	if duplicate := config.DB.
		//Where("name = ?", name).
		Where("name = ? OR code = ?", name, codeNo).
		First(&duplicateLearningArea).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Learning area or Code already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveLearningArea := models2.LearningArea{
		Number:     number,
		Name:       name,
		Code:       codeNo,
		UserNumber: userNumber,
	}
	config.DB.Save(&saveLearningArea)
	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Learning Area Added Successfully!", "data": saveLearningArea})
}

func FetchLearningArea(context *gin.Context) {

	learningAreaResp, err := curriculum.FetchLearningAreaCurriculum("")
	if err != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	var learningAreas_ []responses2.LearningAreaInfoData

	for _, lst := range learningAreaResp {
		learningAreas_ = append(learningAreas_, responses2.LearningAreaInfoData{
			Number:     lst.Number,
			Name:       lst.Name,
			Code:       lst.Code,
			UserNumber: lst.UserNumber,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": learningAreas_})
}
