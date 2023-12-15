package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
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
