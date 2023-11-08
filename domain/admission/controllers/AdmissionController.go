package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewStudentDetails(context *gin.Context) {

	type StudentsDetailsInfo struct {
		Surname            string `json:"surname"`
		SecondName         string `json:"second_name"`
		OtherName          string `json:"other_name"`
		DateOfBirth        string `json:"date_of_birth"`
		BirthCertificateNo string `json:"birth_certificate_no"`
		CountryOfBirth     string `json:"country_of_birth"`
		HomeDistrict       string `json:"home_district"`
		HomeSubCounty      string `json:"home_sub_county"`
		HomeVillage        string `json:"home_village"`
		HomePhone1         string `json:"home_phone1"`
		HomeEmail1         string `json:"home_email1"`
		Religion           string `json:"religion"`
		HomeCountry        string `json:"home_country"`
		HomeCounty         string `json:"home_county"`
		HomeParish         string `json:"home_parish"`
		NearestTown        string `json:"nearest_town"`
		HomePhone2         string `json:"home_phone2"`
		HomeEmail2         string `json:"home_email2"`
	}

	var studentsDetailsInfo StudentsDetailsInfo
	err := context.BindJSON(&studentsDetailsInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Internal sever error!"})
		return
	}

	surname := studentsDetailsInfo.Surname
	secondName := studentsDetailsInfo.SecondName
	otherName := studentsDetailsInfo.OtherName
	dateOfBirth := studentsDetailsInfo.DateOfBirth
	birthCertificateNo := studentsDetailsInfo.BirthCertificateNo
	countryOfBirth := studentsDetailsInfo.CountryOfBirth
	homeDistrict := studentsDetailsInfo.HomeDistrict
	homeSubCounty := studentsDetailsInfo.HomeSubCounty
	homeVillage := studentsDetailsInfo.HomeVillage
	homePhone1 := studentsDetailsInfo.HomePhone1
	homeEmail1 := studentsDetailsInfo.HomeEmail1
	religion := studentsDetailsInfo.Religion
	homeCountry := studentsDetailsInfo.HomeCountry
	homeCounty := studentsDetailsInfo.HomeCounty
	homeParish := studentsDetailsInfo.HomeParish
	nearestTown := studentsDetailsInfo.NearestTown
	homePhone2 := studentsDetailsInfo.HomePhone2
	homeEmail2 := studentsDetailsInfo.HomeEmail2

	userNumber := context.Request.Header["User-Number"][0]

	if surname == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Surname is required!"})
		return
	}

	if secondName == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Second name is required!"})
		return
	}

	if otherName == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Other name is required!"})
		return
	}

	if dateOfBirth == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Date Of Birth is required!"})
		return
	}

	if birthCertificateNo == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Birth certificate is required!"})
		return
	}

	if countryOfBirth == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Country of birth is required!"})
		return
	}

	if homeDistrict == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Home District is required!"})
		return
	}

	if homeSubCounty == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Home SubCounty is required!"})
		return
	}

	if homeVillage == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Home village is required!"})
		return
	}

	if homePhone1 == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Home Phone 1 is required!"})
		return
	}

	if homeEmail1 == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Home email 1 is required!"})
		return
	}

	if religion == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Religion is required!"})
		return
	}

	if homeCountry == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Home Country is required!"})
		return
	}

	if homeCounty == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Home County is required!"})
		return
	}

	if homeParish == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Home Parish is required!"})
		return
	}

	if nearestTown == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Nearest Town is required!"})
		return
	}

	if homePhone2 == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Home Phone 2 is required!"})
		return
	}

	if homeEmail2 == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Home email 2 is required!"})
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

	/*var duplicateCampus campusModel.Campus
	if duplicate := config.DB.
		Where("name = ?", name).
		Where("school_number = ?", schoolNumber).
		First(&duplicateCampus).Error; duplicate == nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Campus already exists"})
		return
	}*/

	/*	number := uuid.NewV4().String()
		saveCampus := studentDetailsModel.StudentDetail{
			Number:       number,
			UserNumber:   userNumber,
			Name:         name,
			Phone:        phone,
			Email:        email,
			Location:     location,
			SchoolNumber: schoolNumber,
		}

		config.DB.Save(&saveCampus)*/

	//context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Campus Added Successfully!", "data": saveCampus})
}
