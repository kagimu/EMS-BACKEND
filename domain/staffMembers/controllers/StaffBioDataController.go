package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	attachmentServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/services"
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

func CreateFamilyBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var request requests.StaffBioDataFamilyRequest
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

	familyNumber := uuid.NewV4().String()

	bioData := models.StaffBioDataFamily{
		Number:                familyNumber,
		JobOfferNumber:        jobOfferNumber,
		FatherSurname:         request.FatherSurname,
		FatherOtherNames:      request.FatherOtherNames,
		FatherEmployment:      request.FatherEmployment,
		FatherPhone1:          request.FatherPhone1,
		FatherPhone2:          request.FatherPhone2,
		FatherEmail1:          request.FatherEmail1,
		FatherEmail2:          request.FatherEmail2,
		FatherPhysicalAddress: request.FatherPhysicalAddress,
		FatherLiving:          request.FatherLiving,
		MotherSurname:         request.MotherSurname,
		MotherOtherNames:      request.MotherOtherNames,
		MotherEmployment:      request.MotherEmployment,
		MotherPhone1:          request.MotherPhone1,
		MotherPhone2:          request.MotherPhone2,
		MotherEmail1:          request.MotherEmail1,
		MotherEmail2:          request.MotherEmail2,
		MotherPhysicalAddress: request.MotherPhysicalAddress,
		MotherLiving:          request.MotherLiving,
		BrotherSurname:        request.BrotherSurname,
		BrotherOtherNames:     request.BrotherOtherNames,
		SisterSurname:         request.SisterSurname,
		SisterOtherNames:      request.SisterOtherNames,
		SpouseSurname:         request.SpouseSurname,
		SpouseMaidenName:      request.SpouseMaidenName,
		SpouseOtherNames:      request.SpouseOtherNames,
		SpouseEmployment:      request.SpouseEmployment,
		SpousePhone1:          request.SpousePhone1,
		SpousePhone2:          request.SpousePhone2,
		SpouseEmail1:          request.SpouseEmail1,
		SpouseEmail2:          request.SpouseEmail2,
		SpousePhysicalAddress: request.SpousePhysicalAddress,
		SpouseLiving:          request.SpouseLiving,
		UserNumber:            ctx.Request.Header.Get("User-Number"),
	}

	var children []models.StaffBioDataFamilyChild
	for _, childReq := range request.Children {
		// Only save if at least one field has data
		if childReq.Surname != "" || childReq.OtherNames != "" || childReq.Dob != "" {
			children = append(children, models.StaffBioDataFamilyChild{
				Number:       uuid.NewV4().String(),
				FamilyNumber: familyNumber,
				Surname:      childReq.Surname,
				OtherNames:   childReq.OtherNames,
				Dob:          childReq.Dob,
			})
		}
	}
	bioData.Children = children

	if err := config.DB.Create(&bioData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Family bio data not saved."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Family bio data saved successfully!", "data": bioData})
}

func FetchFamilyBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var bioData models.StaffBioDataFamily
	if err := config.DB.Preload("Children").Where("job_offer_number = ?", jobOfferNumber).First(&bioData).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Family bio data not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": bioData})
}

func EditFamilyBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var request requests.StaffBioDataFamilyRequest
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

	var bioData models.StaffBioDataFamily
	if err := config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&bioData).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Family bio data not found"})
		return
	}

	bioData.FatherSurname = request.FatherSurname
	bioData.FatherOtherNames = request.FatherOtherNames
	bioData.FatherEmployment = request.FatherEmployment
	bioData.FatherPhone1 = request.FatherPhone1
	bioData.FatherPhone2 = request.FatherPhone2
	bioData.FatherEmail1 = request.FatherEmail1
	bioData.FatherEmail2 = request.FatherEmail2
	bioData.FatherPhysicalAddress = request.FatherPhysicalAddress
	bioData.FatherLiving = request.FatherLiving
	bioData.MotherSurname = request.MotherSurname
	bioData.MotherOtherNames = request.MotherOtherNames
	bioData.MotherEmployment = request.MotherEmployment
	bioData.MotherPhone1 = request.MotherPhone1
	bioData.MotherPhone2 = request.MotherPhone2
	bioData.MotherEmail1 = request.MotherEmail1
	bioData.MotherEmail2 = request.MotherEmail2
	bioData.MotherPhysicalAddress = request.MotherPhysicalAddress
	bioData.MotherLiving = request.MotherLiving
	bioData.BrotherSurname = request.BrotherSurname
	bioData.BrotherOtherNames = request.BrotherOtherNames
	bioData.SisterSurname = request.SisterSurname
	bioData.SisterOtherNames = request.SisterOtherNames
	bioData.SpouseSurname = request.SpouseSurname
	bioData.SpouseMaidenName = request.SpouseMaidenName
	bioData.SpouseOtherNames = request.SpouseOtherNames
	bioData.SpouseEmployment = request.SpouseEmployment
	bioData.SpousePhone1 = request.SpousePhone1
	bioData.SpousePhone2 = request.SpousePhone2
	bioData.SpouseEmail1 = request.SpouseEmail1
	bioData.SpouseEmail2 = request.SpouseEmail2
	bioData.SpousePhysicalAddress = request.SpousePhysicalAddress
	bioData.SpouseLiving = request.SpouseLiving

	// Delete old children
	config.DB.Where("family_number = ?", bioData.Number).Delete(&models.StaffBioDataFamilyChild{})

	// Create new children
	var children []models.StaffBioDataFamilyChild
	for _, childReq := range request.Children {
		if childReq.Surname != "" || childReq.OtherNames != "" || childReq.Dob != "" {
			children = append(children, models.StaffBioDataFamilyChild{
				Number:       uuid.NewV4().String(),
				FamilyNumber: bioData.Number,
				Surname:      childReq.Surname,
				OtherNames:   childReq.OtherNames,
				Dob:          childReq.Dob,
			})
		}
	}
	
	if len(children) > 0 {
		config.DB.Create(&children)
	}

	if err := config.DB.Save(&bioData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Family bio data not updated."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Family bio data updated successfully!"})
}

func CreateEducationBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var request requests.StaffBioDataEducationRequest
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

	educationNumber := uuid.NewV4().String()

	bioData := models.StaffBioDataEducation{
		Number:         educationNumber,
		JobOfferNumber: jobOfferNumber,
		UserNumber:     ctx.Request.Header.Get("User-Number"),
	}

	var items []models.StaffBioDataEducationItem
	
	// Primary
	for _, reqItem := range request.PrimaryEducation {
		if reqItem.Institution != "" || reqItem.Award != "" {
			items = append(items, models.StaffBioDataEducationItem{
				Number:          uuid.NewV4().String(),
				EducationNumber: educationNumber,
				Level:           "Primary",
				Institution:     reqItem.Institution,
				StartYear:       reqItem.StartYear,
				EndYear:         reqItem.EndYear,
				Award:           reqItem.Award,
			})
		}
	}

	// Secondary
	for _, reqItem := range request.SecondaryEducation {
		if reqItem.Institution != "" || reqItem.Award != "" {
			items = append(items, models.StaffBioDataEducationItem{
				Number:          uuid.NewV4().String(),
				EducationNumber: educationNumber,
				Level:           "Secondary",
				Institution:     reqItem.Institution,
				StartYear:       reqItem.StartYear,
				EndYear:         reqItem.EndYear,
				Award:           reqItem.Award,
			})
		}
	}

	// HighSchool
	for _, reqItem := range request.HighSchoolEducation {
		if reqItem.Institution != "" || reqItem.Award != "" {
			items = append(items, models.StaffBioDataEducationItem{
				Number:          uuid.NewV4().String(),
				EducationNumber: educationNumber,
				Level:           "HighSchool",
				Institution:     reqItem.Institution,
				StartYear:       reqItem.StartYear,
				EndYear:         reqItem.EndYear,
				Award:           reqItem.Award,
			})
		}
	}

	// PostHighSchool
	for _, reqItem := range request.PostHighSchoolEducation {
		if reqItem.Institution != "" || reqItem.Award != "" {
			items = append(items, models.StaffBioDataEducationItem{
				Number:          uuid.NewV4().String(),
				EducationNumber: educationNumber,
				Level:           "PostHighSchool",
				Institution:     reqItem.Institution,
				StartDate:       reqItem.StartDate,
				EndDate:         reqItem.EndDate,
				AwardDate:       reqItem.AwardDate,
				Award:           reqItem.Award,
			})
		}
	}

	bioData.Items = items

	if err := config.DB.Create(&bioData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Education bio data not saved."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Education bio data saved successfully!", "data": bioData})
}

func FetchEducationBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var bioData models.StaffBioDataEducation
	if err := config.DB.Preload("Items").Where("job_offer_number = ?", jobOfferNumber).First(&bioData).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Education bio data not found"})
		return
	}

	// Map it back to the request structure for the frontend
	var primary []requests.EducationItemBasic
	var secondary []requests.EducationItemBasic
	var highSchool []requests.EducationItemBasic
	var postHighSchool []requests.EducationItemAdvanced

	for _, item := range bioData.Items {
		switch item.Level {
		case "Primary":
			primary = append(primary, requests.EducationItemBasic{
				Institution: item.Institution,
				StartYear:   item.StartYear,
				EndYear:     item.EndYear,
				Award:       item.Award,
			})
		case "Secondary":
			secondary = append(secondary, requests.EducationItemBasic{
				Institution: item.Institution,
				StartYear:   item.StartYear,
				EndYear:     item.EndYear,
				Award:       item.Award,
			})
		case "HighSchool":
			highSchool = append(highSchool, requests.EducationItemBasic{
				Institution: item.Institution,
				StartYear:   item.StartYear,
				EndYear:     item.EndYear,
				Award:       item.Award,
			})
		case "PostHighSchool":
			postHighSchool = append(postHighSchool, requests.EducationItemAdvanced{
				Institution: item.Institution,
				StartDate:   item.StartDate,
				EndDate:     item.EndDate,
				AwardDate:   item.AwardDate,
				Award:       item.Award,
			})
		}
	}

	response := requests.StaffBioDataEducationRequest{
		PrimaryEducation:        primary,
		SecondaryEducation:      secondary,
		HighSchoolEducation:     highSchool,
		PostHighSchoolEducation: postHighSchool,
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func EditEducationBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var request requests.StaffBioDataEducationRequest
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

	var bioData models.StaffBioDataEducation
	if err := config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&bioData).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Education bio data not found"})
		return
	}

	// Delete old items
	config.DB.Where("education_number = ?", bioData.Number).Delete(&models.StaffBioDataEducationItem{})

	// Create new items
	var items []models.StaffBioDataEducationItem
	
	for _, reqItem := range request.PrimaryEducation {
		if reqItem.Institution != "" || reqItem.Award != "" {
			items = append(items, models.StaffBioDataEducationItem{
				Number:          uuid.NewV4().String(),
				EducationNumber: bioData.Number,
				Level:           "Primary",
				Institution:     reqItem.Institution,
				StartYear:       reqItem.StartYear,
				EndYear:         reqItem.EndYear,
				Award:           reqItem.Award,
			})
		}
	}

	for _, reqItem := range request.SecondaryEducation {
		if reqItem.Institution != "" || reqItem.Award != "" {
			items = append(items, models.StaffBioDataEducationItem{
				Number:          uuid.NewV4().String(),
				EducationNumber: bioData.Number,
				Level:           "Secondary",
				Institution:     reqItem.Institution,
				StartYear:       reqItem.StartYear,
				EndYear:         reqItem.EndYear,
				Award:           reqItem.Award,
			})
		}
	}

	for _, reqItem := range request.HighSchoolEducation {
		if reqItem.Institution != "" || reqItem.Award != "" {
			items = append(items, models.StaffBioDataEducationItem{
				Number:          uuid.NewV4().String(),
				EducationNumber: bioData.Number,
				Level:           "HighSchool",
				Institution:     reqItem.Institution,
				StartYear:       reqItem.StartYear,
				EndYear:         reqItem.EndYear,
				Award:           reqItem.Award,
			})
		}
	}

	for _, reqItem := range request.PostHighSchoolEducation {
		if reqItem.Institution != "" || reqItem.Award != "" {
			items = append(items, models.StaffBioDataEducationItem{
				Number:          uuid.NewV4().String(),
				EducationNumber: bioData.Number,
				Level:           "PostHighSchool",
				Institution:     reqItem.Institution,
				StartDate:       reqItem.StartDate,
				EndDate:         reqItem.EndDate,
				AwardDate:       reqItem.AwardDate,
				Award:           reqItem.Award,
			})
		}
	}

	if len(items) > 0 {
		config.DB.Create(&items)
	}

	if err := config.DB.Save(&bioData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Education bio data not updated."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Education bio data updated successfully!"})
}

func CreateLanguageBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var request requests.StaffBioDataLanguageRequest
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

	recordNumber := uuid.NewV4().String()

	bioData := models.StaffBioDataLanguage{
		Number:         recordNumber,
		JobOfferNumber: jobOfferNumber,
		UserNumber:     ctx.Request.Header.Get("User-Number"),
	}

	var items []models.StaffBioDataLanguageItem
	for _, reqItem := range request.Languages {
		if reqItem.Language != "" {
			items = append(items, models.StaffBioDataLanguageItem{
				Number:               uuid.NewV4().String(),
				LanguageRecordNumber: recordNumber,
				Language:             reqItem.Language,
				Speaking:             reqItem.Speaking,
				Reading:              reqItem.Reading,
				Writing:              reqItem.Writing,
			})
		}
	}
	bioData.Items = items

	if err := config.DB.Create(&bioData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Language bio data not saved."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Language bio data saved successfully!", "data": bioData})
}

func FetchLanguageBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var bioData models.StaffBioDataLanguage
	if err := config.DB.Preload("Items").Where("job_offer_number = ?", jobOfferNumber).First(&bioData).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Language bio data not found"})
		return
	}

	var items []requests.LanguageItemRequest
	for _, item := range bioData.Items {
		items = append(items, requests.LanguageItemRequest{
			Language: item.Language,
			Speaking: item.Speaking,
			Reading:  item.Reading,
			Writing:  item.Writing,
		})
	}

	response := requests.StaffBioDataLanguageRequest{
		Languages: items,
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func EditLanguageBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var request requests.StaffBioDataLanguageRequest
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

	var bioData models.StaffBioDataLanguage
	if err := config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&bioData).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Language bio data not found"})
		return
	}

	// Delete old items
	config.DB.Where("language_record_number = ?", bioData.Number).Delete(&models.StaffBioDataLanguageItem{})

	// Create new items
	var items []models.StaffBioDataLanguageItem
	for _, reqItem := range request.Languages {
		if reqItem.Language != "" {
			items = append(items, models.StaffBioDataLanguageItem{
				Number:               uuid.NewV4().String(),
				LanguageRecordNumber: bioData.Number,
				Language:             reqItem.Language,
				Speaking:             reqItem.Speaking,
				Reading:              reqItem.Reading,
				Writing:              reqItem.Writing,
			})
		}
	}

	if len(items) > 0 {
		config.DB.Create(&items)
	}

	if err := config.DB.Save(&bioData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Language bio data not updated."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Language bio data updated successfully!"})
}

func CreateReferencesBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var request requests.StaffBioDataReferencesRequest
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

	bioData := models.StaffBioDataReferences{
		Number:                 uuid.NewV4().String(),
		JobOfferNumber:         jobOfferNumber,
		Guarantor1Name:         request.Guarantor1.Name,
		Guarantor1Address:      request.Guarantor1.Address,
		Guarantor1Email:        request.Guarantor1.Email,
		Guarantor1Phone:        request.Guarantor1.Phone,
		Guarantor1Relationship: request.Guarantor1.Relationship,
		Guarantor2Name:         request.Guarantor2.Name,
		Guarantor2Address:      request.Guarantor2.Address,
		Guarantor2Email:        request.Guarantor2.Email,
		Guarantor2Phone:        request.Guarantor2.Phone,
		Guarantor2Relationship: request.Guarantor2.Relationship,
		NextOfKinName:          request.NextOfKin.Name,
		NextOfKinAddress:       request.NextOfKin.Address,
		NextOfKinEmail:         request.NextOfKin.Email,
		NextOfKinPhone:         request.NextOfKin.Phone,
		Referee1Name:           request.Referee1.Name,
		Referee1Address:        request.Referee1.Address,
		Referee1Email:          request.Referee1.Email,
		Referee1Phone:          request.Referee1.Phone,
		Referee1Relationship:   request.Referee1.Relationship,
		Referee1YearKnown:      request.Referee1.YearKnown,
		Referee2Name:           request.Referee2.Name,
		Referee2Address:        request.Referee2.Address,
		Referee2Email:          request.Referee2.Email,
		Referee2Phone:          request.Referee2.Phone,
		Referee2Relationship:   request.Referee2.Relationship,
		Referee2YearKnown:      request.Referee2.YearKnown,
		UserNumber:             ctx.Request.Header.Get("User-Number"),
	}

	if err := config.DB.Create(&bioData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. References bio data not saved."})
		return
	}

	// Update Job Offer status since this is the final step
	jobOffer.Status = "Bio Data Completed"
	config.DB.Save(&jobOffer)

	ctx.JSON(http.StatusOK, gin.H{"message": "References bio data saved successfully!", "data": bioData})
}

func FetchReferencesBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var bioData models.StaffBioDataReferences
	if err := config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&bioData).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "References bio data not found"})
		return
	}

	response := requests.StaffBioDataReferencesRequest{
		Guarantor1: requests.GuarantorRequest{
			Name:         bioData.Guarantor1Name,
			Address:      bioData.Guarantor1Address,
			Email:        bioData.Guarantor1Email,
			Phone:        bioData.Guarantor1Phone,
			Relationship: bioData.Guarantor1Relationship,
		},
		Guarantor2: requests.GuarantorRequest{
			Name:         bioData.Guarantor2Name,
			Address:      bioData.Guarantor2Address,
			Email:        bioData.Guarantor2Email,
			Phone:        bioData.Guarantor2Phone,
			Relationship: bioData.Guarantor2Relationship,
		},
		NextOfKin: requests.NextOfKinRequest{
			Name:    bioData.NextOfKinName,
			Address: bioData.NextOfKinAddress,
			Email:   bioData.NextOfKinEmail,
			Phone:   bioData.NextOfKinPhone,
		},
		Referee1: requests.RefereeRequest{
			Name:         bioData.Referee1Name,
			Address:      bioData.Referee1Address,
			Email:        bioData.Referee1Email,
			Phone:        bioData.Referee1Phone,
			Relationship: bioData.Referee1Relationship,
			YearKnown:    bioData.Referee1YearKnown,
		},
		Referee2: requests.RefereeRequest{
			Name:         bioData.Referee2Name,
			Address:      bioData.Referee2Address,
			Email:        bioData.Referee2Email,
			Phone:        bioData.Referee2Phone,
			Relationship: bioData.Referee2Relationship,
			YearKnown:    bioData.Referee2YearKnown,
		},
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func EditReferencesBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	var request requests.StaffBioDataReferencesRequest
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

	var bioData models.StaffBioDataReferences
	if err := config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&bioData).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "References bio data not found"})
		return
	}

	bioData.Guarantor1Name = request.Guarantor1.Name
	bioData.Guarantor1Address = request.Guarantor1.Address
	bioData.Guarantor1Email = request.Guarantor1.Email
	bioData.Guarantor1Phone = request.Guarantor1.Phone
	bioData.Guarantor1Relationship = request.Guarantor1.Relationship

	bioData.Guarantor2Name = request.Guarantor2.Name
	bioData.Guarantor2Address = request.Guarantor2.Address
	bioData.Guarantor2Email = request.Guarantor2.Email
	bioData.Guarantor2Phone = request.Guarantor2.Phone
	bioData.Guarantor2Relationship = request.Guarantor2.Relationship

	bioData.NextOfKinName = request.NextOfKin.Name
	bioData.NextOfKinAddress = request.NextOfKin.Address
	bioData.NextOfKinEmail = request.NextOfKin.Email
	bioData.NextOfKinPhone = request.NextOfKin.Phone

	bioData.Referee1Name = request.Referee1.Name
	bioData.Referee1Address = request.Referee1.Address
	bioData.Referee1Email = request.Referee1.Email
	bioData.Referee1Phone = request.Referee1.Phone
	bioData.Referee1Relationship = request.Referee1.Relationship
	bioData.Referee1YearKnown = request.Referee1.YearKnown

	bioData.Referee2Name = request.Referee2.Name
	bioData.Referee2Address = request.Referee2.Address
	bioData.Referee2Email = request.Referee2.Email
	bioData.Referee2Phone = request.Referee2.Phone
	bioData.Referee2Relationship = request.Referee2.Relationship
	bioData.Referee2YearKnown = request.Referee2.YearKnown

	if err := config.DB.Save(&bioData).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. References bio data not updated."})
		return
	}

	var jobOffer models.JobOffer
	if err := config.DB.Where("number = ?", jobOfferNumber).First(&jobOffer).Error; err == nil {
		jobOffer.Status = "Bio Data Completed"
		config.DB.Save(&jobOffer)
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "References bio data updated successfully!"})
}

func FetchFullStaffProfileBioData(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	response := gin.H{
		"staff_id_number": nil,
		"personal":        nil,
		"contact":         nil,
		"medical":         nil,
		"family":          nil,
		"education":       nil,
		"languages":       nil,
		"references":      nil,
	}

	var deployment models.StaffDeployment
	if err := config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&deployment).Error; err == nil {
		response["staff_id_number"] = deployment.StaffIdNumber
	}

	var personal models.StaffBioDataPersonal
	if err := config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&personal).Error; err == nil {
		response["personal"] = personal
	}

	var contact models.StaffBioDataContact
	if err := config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&contact).Error; err == nil {
		response["contact"] = contact
	}

	var medical models.StaffBioDataMedical
	if err := config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&medical).Error; err == nil {
		response["medical"] = medical
	}

	var family models.StaffBioDataFamily
	if err := config.DB.Preload("Children").Where("job_offer_number = ?", jobOfferNumber).First(&family).Error; err == nil {
		response["family"] = family
	}

	var education models.StaffBioDataEducation
	if err := config.DB.Preload("Items").Where("job_offer_number = ?", jobOfferNumber).First(&education).Error; err == nil {
		var primary []requests.EducationItemBasic
		var secondary []requests.EducationItemBasic
		var highSchool []requests.EducationItemBasic
		var postHighSchool []requests.EducationItemAdvanced

		for _, item := range education.Items {
			switch item.Level {
			case "Primary":
				primary = append(primary, requests.EducationItemBasic{Institution: item.Institution, StartYear: item.StartYear, EndYear: item.EndYear, Award: item.Award})
			case "Secondary":
				secondary = append(secondary, requests.EducationItemBasic{Institution: item.Institution, StartYear: item.StartYear, EndYear: item.EndYear, Award: item.Award})
			case "HighSchool":
				highSchool = append(highSchool, requests.EducationItemBasic{Institution: item.Institution, StartYear: item.StartYear, EndYear: item.EndYear, Award: item.Award})
			case "PostHighSchool":
				postHighSchool = append(postHighSchool, requests.EducationItemAdvanced{Institution: item.Institution, StartDate: item.StartDate, EndDate: item.EndDate, AwardDate: item.AwardDate, Award: item.Award})
			}
		}
		response["education"] = requests.StaffBioDataEducationRequest{
			PrimaryEducation: primary, SecondaryEducation: secondary, HighSchoolEducation: highSchool, PostHighSchoolEducation: postHighSchool,
		}
	}

	var language models.StaffBioDataLanguage
	if err := config.DB.Preload("Items").Where("job_offer_number = ?", jobOfferNumber).First(&language).Error; err == nil {
		var items []requests.LanguageItemRequest
		for _, item := range language.Items {
			items = append(items, requests.LanguageItemRequest{Language: item.Language, Speaking: item.Speaking, Reading: item.Reading, Writing: item.Writing})
		}
		response["languages"] = requests.StaffBioDataLanguageRequest{Languages: items}
	}

	var references models.StaffBioDataReferences
	if err := config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&references).Error; err == nil {
		response["references"] = requests.StaffBioDataReferencesRequest{
			Guarantor1: requests.GuarantorRequest{Name: references.Guarantor1Name, Address: references.Guarantor1Address, Email: references.Guarantor1Email, Phone: references.Guarantor1Phone, Relationship: references.Guarantor1Relationship},
			Guarantor2: requests.GuarantorRequest{Name: references.Guarantor2Name, Address: references.Guarantor2Address, Email: references.Guarantor2Email, Phone: references.Guarantor2Phone, Relationship: references.Guarantor2Relationship},
			NextOfKin:  requests.NextOfKinRequest{Name: references.NextOfKinName, Address: references.NextOfKinAddress, Email: references.NextOfKinEmail, Phone: references.NextOfKinPhone},
			Referee1:   requests.RefereeRequest{Name: references.Referee1Name, Address: references.Referee1Address, Email: references.Referee1Email, Phone: references.Referee1Phone, Relationship: references.Referee1Relationship, YearKnown: references.Referee1YearKnown},
			Referee2:   requests.RefereeRequest{Name: references.Referee2Name, Address: references.Referee2Address, Email: references.Referee2Email, Phone: references.Referee2Phone, Relationship: references.Referee2Relationship, YearKnown: references.Referee2YearKnown},
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func UploadStaffPhoto(ctx *gin.Context) {
	jobOfferNumber := ctx.Param("job_offer_number")
	if jobOfferNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Job offer number is required"})
		return
	}

	attachmentNumber, err := attachmentServices.UploadAttachment(ctx, "file")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload photo: " + err.Error()})
		return
	}

	var personal models.StaffBioDataPersonal
	err = config.DB.Where("job_offer_number = ?", jobOfferNumber).First(&personal).Error
	if err != nil {
		personal = models.StaffBioDataPersonal{
			Number:         uuid.NewV4().String(),
			JobOfferNumber: jobOfferNumber,
			Photo:          attachmentNumber,
		}
		if err := config.DB.Create(&personal).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save photo to personal profile"})
			return
		}
	} else {
		personal.Photo = attachmentNumber
		if err := config.DB.Save(&personal).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update photo in personal profile"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Photo uploaded successfully", "data": personal, "photo": attachmentNumber})
}

