package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/learningMaterials/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/learningMaterials/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
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
	learningMaterial.Category = ctx.PostForm("category")
	learningMaterial.SubCategoryNumber = ctx.PostForm("sub_category")
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

	if err := config.DB.Create(&learningMaterial).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "An error occurred. Learning material not saved"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Learning material created successfully!"})

}

func DeleteMaterial(ctx *gin.Context) {

	number := ctx.Param("number")
	var learningMaterial models.LearningMaterial

	err := config.DB.Where(&models.LearningMaterial{Number: number}).First(&learningMaterial).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Learning material not found"})
	}

	err = config.DB.Delete(&models.LearningMaterial{}, "number = ?", number).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Learning material not deleted try again " + err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Learning material deleted successfully!"})

}

func UpdateLearningMaterial(ctx *gin.Context) {
	var learningMaterial = models.LearningMaterial{}
	/*	err := ctx.BindJSON(&learningMaterial)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
			return
		}*/

	learningMaterial.Name = ctx.PostForm("name")
	learningMaterial.Type = ctx.PostForm("type")
	learningMaterial.Category = ctx.PostForm("category")
	learningMaterial.SubCategoryNumber = ctx.PostForm("sub_category")
	learningMaterialNumber := ctx.PostForm("learning_material_number")

	if learningMaterialNumber != "" {
		learningMaterial.Number = learningMaterialNumber
	}

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

	if file != nil {
		// Upload the file to specific dst.
		err = ctx.SaveUploadedFile(file, "fileRepo/learningMaterials/"+file.Filename)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return

		}

	}

	if err := config.DB.Where("number", learningMaterialNumber).
		Updates(learningMaterial).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "An error occurred. Learning material not saved " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Learning material updated successfully!"})

}

func FetchAllLearningMaterials(ctx *gin.Context) {
	var learningMaterials []models.LearningMaterial
	config.DB.Order("code asc").Find(&learningMaterials)

	ctx.JSON(http.StatusOK, gin.H{"data": learningMaterials})

}

func CreateCurriculumReference(context *gin.Context) {

	type References struct {
		Number                 string `json:"number"`
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
		//if reference has a number we just update else we submit a new one
		if refs.Number != "" {

			var referenceDataInfo models.CurriculumSubtopicReference
			updateRef := tx.Model(&referenceDataInfo).
				Where("number = ?", refs.Number).
				Updates(models.CurriculumSubtopicReference{Description: refs.Description, LearningMaterialNumber: refs.LearningMaterialNumber, Code: refs.Code})
			if updateRef.Error != nil {
				context.JSON(http.StatusNotFound, gin.H{"message": "Failed to update Reference!"})
				return
			}

		} else {

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
	}

	// Commit the transaction if there is no error
	if err := tx.Commit().Error; err != nil {
		// Handle or return the error as needed
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create reference!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Success"})
}

func CreateMaterialSubType(context *gin.Context) {
	var materialSubType = models.MaterialSubType{}
	var request = requests.MaterialSubTypeRequest{}
	err := context.BindJSON(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Please structure the response according to specifications!"})
		return
	}

	err = request.Validate()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	materialSubType.UserNumber = context.Request.Header["User-Number"][0]
	materialSubType.Name = request.Name
	materialSubType.Category = request.Category
	materialSubType.Number = uuid.NewV4().String()
	// convert string to int64
	materialSubType.Code, _ = strconv.ParseInt(request.Code, 10, 64)

	config.DB.Create(&materialSubType)

	context.JSON(http.StatusOK, gin.H{"message": "Material sub type created successfully!"})

}

func FetchMaterialSubTypesByCategory(context *gin.Context) {
	var category = context.Param("category")

	var materialSubTypes []models.MaterialSubType
	config.DB.Where("category = ? ", category).Find(&materialSubTypes)

	context.JSON(http.StatusOK, gin.H{"data": materialSubTypes})

}
