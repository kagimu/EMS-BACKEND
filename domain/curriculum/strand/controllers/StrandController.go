package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum"
	responses4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/api/responses"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/api/responses"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/models"
	responses3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/api/responses"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func NewStrand(context *gin.Context) {

	type SubjectClusterInfo struct {
		Code          string `json:"code"`
		Name          string `json:"name"`
		SubjectNumber string `json:"subject_number"`
	}

	var subjectClusterInfo SubjectClusterInfo
	err := context.BindJSON(&subjectClusterInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	code := subjectClusterInfo.Code
	name := subjectClusterInfo.Name
	subjectNumber := subjectClusterInfo.SubjectNumber

	userNumber := context.Request.Header["User-Number"][0]

	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Code is required!"})
		return
	}

	if name == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Name is required!"})
		return
	}

	if subjectNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Subject is required!"})
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

	var duplicateStrand models2.Strand
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("subject_number = ?", subjectNumber).
		First(&duplicateStrand).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Strand already exists"})
		return
	}

	number := uuid.NewV4().String()
	///save user data
	saveStrand := models2.Strand{
		Number:        number,
		Name:          name,
		SubjectNumber: subjectNumber,
		UserNumber:    userNumber,
		Code:          codeNo,
	}
	config.DB.Save(&saveStrand)

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Subject cluster Added Successfully!", "data": saveStrand})
}

func FetchStrand(context *gin.Context) {

	//learningAreaNumber := context.Query("learning_area_number")
	//subjectClusterNumber := context.Query("subject_cluster_number")
	subjectNumber := context.Query("subject_number")

	//fetch learning area
	learningAreaResp, err := curriculum.FetchLearningAreaCurriculum("")
	if err != "" {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
		return
	}

	var learningAreas_ []responses4.LearningAreaInfoData

	for _, lstLA := range learningAreaResp {

		//fetch subject cluster
		subjectClusterResp, err := curriculum.FetchSubjectClusterCurriculum(lstLA.Number, "")
		if err != "" {
			context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
			return
		}

		var subjectCluster_ []responses3.SubjectClusterInfo
		for _, lstSC := range subjectClusterResp {

			subjectResp, err := curriculum.FetchSubjectCurriculum(lstSC.Number, subjectNumber)
			if err != "" {
				context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
				return
			}

			var subjects_ []responses.SubjectInfoList
			for _, lstSB := range subjectResp {

				strandResp, err := curriculum.FetchStrandCurriculum(lstSB.Number, "")
				if err != "" {
					context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err})
					return
				}

				var strand_ []responses2.StrandDataList
				for _, lstST := range strandResp {
					strand_ = append(strand_, responses2.StrandDataList{
						Number:        lstST.Number,
						Name:          lstST.Name,
						Code:          lstST.Code,
						UserNumber:    lstST.UserNumber,
						SubjectNumber: lstST.SubjectNumber,
					})
				}

				subjects_ = append(subjects_, responses.SubjectInfoList{
					Number:               lstSB.Number,
					Name:                 lstSB.Name,
					Code:                 lstSB.Code,
					UserNumber:           lstSB.UserNumber,
					SubjectClusterNumber: lstSB.SubjectClusterNumber,
					Strand:               strand_,
				})

			}

			subjectCluster_ = append(subjectCluster_, responses3.SubjectClusterInfo{
				Number:             lstSC.Number,
				Name:               lstSC.Name,
				Code:               lstSC.Code,
				LearningAreaNumber: lstSC.LearningAreaNumber,
				Subject:            subjects_,
			})

		}

		learningAreas_ = append(learningAreas_, responses4.LearningAreaInfoData{
			Number:         lstLA.Number,
			Name:           lstLA.Name,
			Code:           lstLA.Code,
			UserNumber:     lstLA.UserNumber,
			SubjectCluster: subjectCluster_,
		})

	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": learningAreas_})

}

func subject(subjectNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if subjectNumber != "" {
			return db.Where("subject_number = ?", subjectNumber)
		}
		return db
	}
}
