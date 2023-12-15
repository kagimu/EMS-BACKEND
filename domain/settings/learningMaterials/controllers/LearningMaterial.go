package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/learningMaterials/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"strconv"
)

func CreateNewLearningMaterial(ctx *gin.Context) {
	var learningMaterial = models.LearningMaterial{}
	/*	err := ctx.BindJSON(&learningMaterial)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
			return
		}*/

	learningMaterial.Number = uuid.NewV4().String()
	learningMaterial.Name = ctx.PostForm("name")
	learningMaterial.Type = ctx.PostForm("type")
	Code, err := strconv.Atoi(ctx.PostForm("code"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}
	learningMaterial.Code = int64(Code)
	learningMaterial.UserNumber = ctx.Request.Header["User-Number"][0]

	/*err := learningMaterial.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}*/

	file, _ := ctx.FormFile("file")
	log.Println(file.Filename)

	// Upload the file to specific dst.
	err = ctx.SaveUploadedFile(file, "fileRepo/learningMaterials/"+file.Filename)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return

	}

	config.DB.Create(&learningMaterial)

	ctx.JSON(http.StatusOK, gin.H{"message": "Learning material created successfully!"})

}

func FetchAllLearningMaterials(ctx *gin.Context) {
	var learningMaterials []models.LearningMaterial
	config.DB.Find(&learningMaterials)

	ctx.JSON(http.StatusOK, gin.H{"data": learningMaterials})

}

func CreateCurriculumReference(context *gin.Context) {

	type ReferenceInfo struct {
		Code             int64  `json:"material_code"`
		MaterialNumber   string `json:"material_number"`
		CurriculumNumber string `json:"curriculum_number"`
		SubTopicNumber   string `json:"sub_topic_number"`
	}

	var referenceInfo ReferenceInfo
	err := context.BindJSON(&referenceInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err})
		//context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	code := referenceInfo.Code
	materialNumber := referenceInfo.MaterialNumber
	curriculumNumber := referenceInfo.CurriculumNumber
	subTopicNumber := referenceInfo.SubTopicNumber
	userNumber := context.Request.Header["User-Number"][0]

	/*if code > 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code is required!"})
		return
	}*/

	if materialNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Resource is required!"})
		return
	}

	if curriculumNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Curriculum is required!"})
		return
	}

	if subTopicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subtopic is required!"})
		return
	}

	var userData models2.User
	//check if user is a valid user
	validateUser := config.DB.
		Where("number = ?", userNumber).
		First(&userData)
	if validateUser.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	/*	codeNo, err := strconv.ParseInt(code, 10, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}*/

	tx := config.DB.Begin()

	///save user data
	saveReference := models.CurriculumSubtopicReference{
		Number:                 uuid.NewV4().String(),
		Code:                   code,
		LearningMaterialNumber: materialNumber,
		CurriculumNumber:       curriculumNumber,
		SubTopicNumber:         subTopicNumber,
		UserNumber:             userNumber,
	}

	if err := tx.Create(&saveReference).Error; err != nil {
		// If any error occurs, rollback the transaction
		tx.Rollback()
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to save reference"})
		return
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create reference!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Reference Added Successfully!", "data": saveReference})
}
