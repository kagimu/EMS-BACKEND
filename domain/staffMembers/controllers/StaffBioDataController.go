package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreatePersonalBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var request requests.StaffBioDataPersonalRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please structure the response according to specifications!"})
		return
	}

	err = request.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Verify Job Offer exists
	var jobOffer models.JobOffer
	if err := config.DB.Where("number = ?", jobOfferNumber).First(&jobOffer).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Job offer not found"})
		return
	}

	bioData := models.StaffBioDataPersonal{
		Number:             uuid.NewV4().String(),
		JobOfferNumber:     jobOfferNumber,
		Surname:            request.Surname,
		SecondName:         request.SecondName,
		OtherNames:         request.OtherNames,
		DateCompleted:      request.DateCompleted,
		DateOfBirth:        request.DateOfBirth,
		Sex:                request.Sex,
		BirthCertificateNo: request.BirthCertificateNo,
		Religion:           request.Religion,
		NationalIdNo:       request.NationalIdNo,
		CountryOfBirth:     request.CountryOfBirth,
		HomeCountry:        request.HomeCountry,
		HomeDistrict:       request.HomeDistrict,
		HomeCounty:         request.HomeCounty,
		HomeSubCounty:      request.HomeSubCounty,
		HomeParish:         request.HomeParish,
		HomeVillage:        request.HomeVillage,
		VillageOfResidence: request.VillageOfResidence,
		TaxIdNo:            request.TaxIdNo,
		UserNumber:         ctx.Request.Header.Get("User-Number"),
	}

	if err := config.DB.Create(&bioData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Bio data not saved."})
		return
	}

	// Update JobOffer status
	jobOffer.Status = "Filling Bio Data"
	config.DB.Save(&jobOffer)

	ctx.JSON(http.StatusOK, gin.H{"message": "Personal bio data saved successfully!", "data": bioData})
}

func FetchPersonalBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var bioData models.StaffBioDataPersonal
	if err := config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&bioData).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Personal bio data not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": bioData})
}

func EditPersonalBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var request requests.StaffBioDataPersonalRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please structure the response according to specifications!"})
		return
	}

	err = request.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var bioData models.StaffBioDataPersonal
	if err := config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&bioData).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Personal bio data not found"})
		return
	}

	bioData.Surname = request.Surname
	bioData.SecondName = request.SecondName
	bioData.OtherNames = request.OtherNames
	bioData.DateCompleted = request.DateCompleted
	bioData.DateOfBirth = request.DateOfBirth
	bioData.Sex = request.Sex
	bioData.BirthCertificateNo = request.BirthCertificateNo
	bioData.Religion = request.Religion
	bioData.NationalIdNo = request.NationalIdNo
	bioData.CountryOfBirth = request.CountryOfBirth
	bioData.HomeCountry = request.HomeCountry
	bioData.HomeDistrict = request.HomeDistrict
	bioData.HomeCounty = request.HomeCounty
	bioData.HomeSubCounty = request.HomeSubCounty
	bioData.HomeParish = request.HomeParish
	bioData.HomeVillage = request.HomeVillage
	bioData.VillageOfResidence = request.VillageOfResidence
	bioData.TaxIdNo = request.TaxIdNo

	if err := config.DB.Save(&bioData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Bio data not updated."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Personal bio data updated successfully!"})
}

func CreateContactBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var request requests.StaffBioDataContactRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please structure the response according to specifications!"})
		return
	}

	err = request.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Verify Job Offer exists
	var jobOffer models.JobOffer
	if err := config.DB.Where("number = ?", jobOfferNumber).First(&jobOffer).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Job offer not found"})
		return
	}

	bioData := models.StaffBioDataContact{
		Number:          uuid.NewV4().String(),
		JobOfferNumber:  jobOfferNumber,
		MainEmail:       request.MainEmail,
		SecondEmail:     request.SecondEmail,
		AltEmail:        request.AltEmail,
		MainPhone:       request.MainPhone,
		SecondPhone:     request.SecondPhone,
		AltPhone:        request.AltPhone,
		PostMailAddress: request.PostMailAddress,
		PhysicalAddress: request.PhysicalAddress,
		UserNumber:      ctx.Request.Header.Get("User-Number"),
	}

	if err := config.DB.Create(&bioData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Contact bio data not saved."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Contact bio data saved successfully!", "data": bioData})
}

func FetchContactBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var bioData models.StaffBioDataContact
	if err := config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&bioData).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Contact bio data not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": bioData})
}

func EditContactBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var request requests.StaffBioDataContactRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please structure the response according to specifications!"})
		return
	}

	err = request.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var bioData models.StaffBioDataContact
	if err := config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&bioData).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Contact bio data not found"})
		return
	}

	bioData.MainEmail = request.MainEmail
	bioData.SecondEmail = request.SecondEmail
	bioData.AltEmail = request.AltEmail
	bioData.MainPhone = request.MainPhone
	bioData.SecondPhone = request.SecondPhone
	bioData.AltPhone = request.AltPhone
	bioData.PostMailAddress = request.PostMailAddress
	bioData.PhysicalAddress = request.PhysicalAddress

	if err := config.DB.Save(&bioData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Contact bio data not updated."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Contact bio data updated successfully!"})
}

func CreateMedicalBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var request requests.StaffBioDataMedicalRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please structure the response according to specifications!"})
		return
	}

	err = request.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var jobOffer models.JobOffer
	if err := config.DB.Where("number = ?", jobOfferNumber).First(&jobOffer).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Job offer not found"})
		return
	}

	bioData := models.StaffBioDataMedical{
		Number:              uuid.NewV4().String(),
		JobOfferNumber:      jobOfferNumber,
		Height:              request.Height,
		Weight:              request.Weight,
		Complexion:          request.Complexion,
		BloodType:           request.BloodType,
		Allergies:           request.Allergies,
		MedicalConditions:   request.MedicalConditions,
		RegularMedications:  request.RegularMedications,
		PhysicalImpairments: request.PhysicalImpairments,
		DoctorName:          request.DoctorName,
		DoctorPhone:         request.DoctorPhone,
		DoctorFacility:      request.DoctorFacility,
		UserNumber:          ctx.Request.Header.Get("User-Number"),
	}

	if err := config.DB.Create(&bioData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Medical bio data not saved."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Medical bio data saved successfully!", "data": bioData})
}

func FetchMedicalBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var bioData models.StaffBioDataMedical
	if err := config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&bioData).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Medical bio data not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": bioData})
}

func EditMedicalBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var request requests.StaffBioDataMedicalRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please structure the response according to specifications!"})
		return
	}

	err = request.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var bioData models.StaffBioDataMedical
	if err := config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&bioData).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Medical bio data not found"})
		return
	}

	bioData.Height = request.Height
	bioData.Weight = request.Weight
	bioData.Complexion = request.Complexion
	bioData.BloodType = request.BloodType
	bioData.Allergies = request.Allergies
	bioData.MedicalConditions = request.MedicalConditions
	bioData.RegularMedications = request.RegularMedications
	bioData.PhysicalImpairments = request.PhysicalImpairments
	bioData.DoctorName = request.DoctorName
	bioData.DoctorPhone = request.DoctorPhone
	bioData.DoctorFacility = request.DoctorFacility

	if err := config.DB.Save(&bioData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Medical bio data not updated."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Medical bio data updated successfully!"})
}
