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

	type References struct {
		Code                   int64  `json:"code"`
		LearningMaterialNumber string `json:"learningMaterialNumber"`
		Description            string `json:"description"`
	}

	type ReferenceInfo struct {
		SubTopicNumber string       `json:"sub_topic_number"`
		References     []References `json:"references"`
	}

	var referenceInfo ReferenceInfo
	err := context.BindJSON(&referenceInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err})
		//context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	subTopicNumber := referenceInfo.SubTopicNumber
	references := referenceInfo.References

	userNumber := context.Request.Header["User-Number"][0]

	if subTopicNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subtopic is required!"})
		return
	}

	if len(references) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "References are required!"})
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

	tx := config.DB.Begin()

	for _, refs := range references {
		///save user data
		saveReference := models.CurriculumSubtopicReference{
			Number:                 uuid.NewV4().String(),
			Code:                   refs.Code,
			LearningMaterialNumber: refs.LearningMaterialNumber,
			Description:            refs.Description,
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
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create reference!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Reference Added Successfully!"})
}
