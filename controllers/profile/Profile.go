package profile

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/config"
	"github.com/hamsoftug/wholesomeEd/models"
	"net/http"
)

func Profile(context *gin.Context) {

	userId := ""
	if len(context.Request.Header["User-Number"]) > 0 {
		userId = context.Request.Header["User-Number"][0]
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found!"})
		return
	}

	var userProfile models.User

	if err := config.DB.
		Where("number = ?", userId).
		Preload("Contacts").
		First(&userProfile).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "User not found"})
		return
	}

	var contactListInfo []models.ContactList

	for _, item := range userProfile.Contacts {
		contactListInfo = append(contactListInfo, models.ContactList{
			ID:          item.ID,
			Contact:     item.Contact,
			ContactType: item.ContactType,
			//UserNumber:  item.UserNumber,
		})
	}

	_userProfile := models.GetUserList{
		ID:          userProfile.ID,
		Number:      userProfile.Number,
		Name:        userProfile.Name,
		DateOfBirth: userProfile.DateOfBirth,
		Gender:      userProfile.Gender,
		CoverPic:    CoverProfilePic(context, "cover_photo", userId, "defaultCoverPic.png"),
		ProfilePic:  CoverProfilePic(context, "profile", userId, "default.png"),
		ContactList: contactListInfo,
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _userProfile})
}

func CoverProfilePic(context *gin.Context, pictureType string, userId string, defaultPic string) string {
	profilePath := ""
	//path := "https://" + context.Request.Host + "/fileRepo/profileCoverPhoto/"
	path := "https://" + context.Request.Host + "/api/v1/main/file/"
	var displayPic models.DisplayPic
	if err := config.DB.
		Where("picture_type = ?", pictureType).
		Where("user_number = ?", userId).
		Order("created_at desc").
		First(&displayPic).Error; err != nil {
		profilePath = path + defaultPic + "?cat=cover"
	} else {
		profilePath = path + displayPic.FileName + "?cat=cover"
	}
	return profilePath
}
