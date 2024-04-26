package commons

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func DeleteFileFromPath(path string, context *gin.Context) {
	//path := "/fileRepo/excelFile/" + pdfDel.FileName
	// detect file if it exists
	_, err2 := mimetype.DetectFile(path)
	if err2 == nil {
		//a file is removed if it exists in the path
		//removing png file
		errRemove := os.Remove(path)
		if errRemove != nil {
			context.JSON(http.StatusBadRequest, gin.H{"message": errRemove})
			return
		}

	}
}

func CheckValidPhone(phone string) string {

	replaceCharacter := strings.NewReplacer("(", "", ")", "", "-", "", " ", "")
	phoneNoChar := replaceCharacter.Replace(phone)

	phone = strings.TrimSpace(strings.ReplaceAll(phoneNoChar, "+", ""))
	_, err := strconv.Atoi(phone)
	if err != nil {
		return ""
	}

	//check if phone starts plus
	phoneStart := strings.HasPrefix(phoneNoChar, "+")
	phoneLength := len(phoneNoChar)

	//fmt.Println("Result 1: ", phoneLength, phoneStart)
	if phoneLength <= 8 || phoneLength >= 14 {
		return ""
	}

	if phoneStart {
		if phoneLength > 10 {
			return phoneNoChar
		} else {
			return ""
		}
	}

	if phoneLength == 10 {
		firstPhoneChar := phoneNoChar[0:1]
		if firstPhoneChar == "0" {
			replaceZero := phoneNoChar[1:]
			return "+256" + replaceZero
		} else {
			return "+" + phoneNoChar
		}

		//fmt.Println("Result 1: ", firstPhoneChar)

	} else if phoneLength == 9 {
		firstPhoneChar := phoneNoChar[0:1]
		if firstPhoneChar == "7" {
			return "+256" + phoneNoChar
		} else {
			return ""
		}

	} else {
		return "+" + phoneNoChar
	}

}

func Feedback(message string) map[string]any {
	response := make(map[string]any)
	response["message"] = message
	return response
}

func GetUserFromHeader(c *gin.Context) (*models.User, error) {
	userNumber := c.Request.Header["User-Number"][0]
	var user models.User

	if err := config.DB.
		Where("number = ?", userNumber).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

type PermissionGroup map[string][]string

func PermissionGroups() []PermissionGroup {

	return []PermissionGroup{
		{"USERS": []string{"create_user", "view_users", "update_user_details", "delete"}},
		{"PERMISSIONS": []string{"create_permission", "view_permission", "delete"}},
		{"CURRICULUM": []string{"create_curriculum", "view_curriculum", "export_curriculum",
			"delete_curriculum", "import_curriculum", "view_subject-cluster", "view_subject", "view_strands", "view_topic", "view_sub-topic", "view_concept"}},
		{"SCHEME OF WORK": []string{"view_scheme-of-work", "create_scheme", "import_scheme", "export_scheme", "delete_scheme"}},
		{"SCHOOL": []string{"view_school", "add_school", "edit_school", "delete_school"}},
		{"TIMETABLE": []string{"add_timetable", "view_timetable-details"}},
		{"STUDENTS": []string{"view_students", "update_student_details"}},
		{"CALENDAR": []string{"view_calendar", "view_calendar-details"}},
		{"INTERVIEWS": []string{"view_interview-details"}},
		{"PRODUCTS & SERVICES": []string{"view_services", "view_products"}},
		{"USER ACCOUNTS": []string{"view_account-management", "view_user-management"}},
		{"ASSESSMENT": []string{"view_grade-scale", "view_grade-standard", "view_assessment-list", "view_examtype"}},
		{"LEARNING AREAS": []string{"create_learning_area", "view_learning_area", "update_learning_area"}},
		{"SETTINGS": []string{"create_learning_area", "view_learning_area", "update_learning_area"}},
	}

}

func HandleFileUpload(c *gin.Context, formFileKey string, savePath string) (string, error) {
	attachmentFile, err := c.FormFile(formFileKey)
	if err != nil {
		return "", err
	}

	if attachmentFile.Size > 20971520 {
		return " ", errors.New("file size should be less than 20MB")
	}

	t := time.Now()
	name := GeneratedFormattedUUID() + "-" + strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()) + "-" + strconv.Itoa(t.Minute())
	extension := filepath.Ext(attachmentFile.Filename)
	if extension != ".pdf" {
		return " ", errors.New("file should be a .pdf")
	}

	fileName := name + extension
	filePath := savePath + fileName

	if err := c.SaveUploadedFile(attachmentFile, filePath); err != nil {
		return "", err
	}

	return fileName, nil
}

func GeneratedFormattedUUID() string {
	u := uuid.NewV4().String()
	formattedUUID := strings.ReplaceAll(u, "-", "")
	return formattedUUID
}

func CheckDay(day string) (string, bool) {
	validDays := map[string]bool{
		"Monday":    true,
		"Tuesday":   true,
		"Wednesday": true,
		"Thursday":  true,
		"Friday":    true,
		"Saturday":  true,
		"Sunday":    true,
	}

	if valid := validDays[day]; valid {
		return "", true
	}

	return "Invalid day", false
}
