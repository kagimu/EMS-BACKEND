package userManagement

import (
	"github.com/gin-gonic/gin"
	"github.com/hamsoftug/wholesomeEd/config"
	"github.com/hamsoftug/wholesomeEd/controllers/profile"
	"github.com/hamsoftug/wholesomeEd/models"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func Users(context *gin.Context) {

	userNumber := ""
	if len(context.Request.Header["User-Number"]) > 0 {
		userNumber = context.Request.Header["User-Number"][0]
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User not found!"})
		return
	}

	//	validate user number
	var user models.User
	if validateUserNumber := config.DB.
		Where("number = ?", userNumber).
		Where("type = ?", "Admin").
		First(&user).Error; validateUserNumber != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User Not found!"})
		return
	}

	if user.Type == "User" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Bad Request!"})
		return
	}

	page := context.Query("page")
	perPage := context.Query("per_page")
	search := context.Query("search")
	userType := context.Query("userType")

	if userType == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid data!"})
		return
	}

	if page == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Page Number is required!"})
		return
	}

	if perPage == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Per page Number is required!"})
		return
	}

	pageIntNumber, err := strconv.Atoi(page)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Failed to convert page to int"})
	}

	perPageOffset, err := strconv.Atoi(perPage)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Failed to convert page to int"})
	}

	pageOffset := (pageIntNumber - 1) * perPageOffset

	var userList []models.User
	var _userList []models.GetUserList

	db := config.DB
	if err := db.
		Scopes(filterUserData(userType)).
		Where(
			db.Where("contacts.contact like ?", "%"+search+"%").
				Or("name like ?", "%"+search+"%")).
		Joins("join contacts on contacts.user_number = users.number").
		Preload("Contacts").
		Group("number").
		Offset(pageOffset).
		Limit(perPageOffset).
		Order("created_at desc").
		Find(&userList).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Profile not found"})
		return
	}

	for _, item := range userList {

		var contactListInfo []models.ContactList
		for _, item := range item.Contacts {
			contactListInfo = append(contactListInfo, models.ContactList{
				ID:          item.ID,
				Contact:     item.Contact,
				ContactType: item.ContactType,
			})
		}

		_userList = append(_userList, models.GetUserList{
			ID:          item.ID,
			Number:      item.Number,
			Name:        item.Name,
			DateOfBirth: item.DateOfBirth,
			Gender:      item.Gender,
			Type:        item.Type,
			Status:      item.Status,
			Info:        item.Info,
			Seen:        item.Seen,
			CreatedAt:   item.CreatedAt,
			CoverPic:    profile.CoverProfilePic(context, "cover_photo", item.Number, "defaultCoverPic.png"),
			ProfilePic:  profile.CoverProfilePic(context, "profile", item.Number, "default.png"),
			ContactList: contactListInfo,
		})
	}

	context.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _userList, "per_page": perPageOffset, "page": pageIntNumber})
}

func filterUserData(userType string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if userType != "" {
			return db.Where("users.type = ?", userType)
		}
		return db
	}
}
