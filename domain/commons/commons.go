package commons

import (
	"errors"
	"math"
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
		{"CURRICULUM": []string{"create_curriculum", "view_curriculum", "export_curriculum",
			"delete_curriculum", "import_curriculum", "create_learning_area", "view_learning_area", "update_learning_area", "view_subject-cluster", "view_subject",
			"view_strands", "view_topic", "view_sub-topic", "view_concept"}},
		{"SCHEME OF WORK": []string{"view_scheme-of-work", "create_scheme", "import_scheme", "export_scheme", "delete_scheme"}},
		{"SETUP": []string{"view_school", "add_school", "edit_school", "delete_school", "create_campus", "edit_campus", "delete_campus", "create_academic-year", "view_academic-year", "update_academic-year",
			"create_term", "view_term", "update_term", "create_class", "view_class", "update_class", "create_stream", "view_streams", "update_stream"}},
		{"USERS": []string{"create_user", "view_users", "update_user_details", "delete_user", "view_account-management", "view_user-management", "add_timetable", "view_timetable-details"}},
		{"PERMISSIONS": []string{"create_permission", "view_permission", "delete_permission", "update_permission"}},

		{"STUDENTS": []string{"view_students", "update_student_details"}},
		{"CALENDAR": []string{"view_calendar", "view_calendar-details"}},
		{"INTERVIEWS": []string{"view_interview-details", "create_interview_schedules"}},
		{"Catering": []string{"view_catering_menu", "create_unit_of_measure",
			"record_meals", "view_meals", "view_unit_of_measure", "create_feeding_groups", ""}},
		{"PRODUCTS & SERVICES": []string{"create_product", "create_service", "edit_product", "edit_service", "view_services", "view_products"}},
		{"ASSESSMENT": []string{"view_grade-scale", "view_grade-standard", "view_assessment-list", "view_examtype"}},
		{"ADMISSION": []string{"admit_student", "view_admissions", "cancel_admission"}},
		{"ENROLLMENT": []string{"enroll_student", "view_enrollments", "terminate_enrollment", "re-instate_enrollment"}},
		{"SETTINGS": []string{"access_settings_menu", "view_learning_area", "update_learning_area"}},
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

func CheckDay(day string) string {

	for _, _day := range DaysOfTheWeek() {

		if day == _day {
			return ""
		}

	}

	return "Invalid day: " + day
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func DaysOfTheWeek() []string {

	now := time.Now()

	// Find the first Monday on or before the current time
	monday := now.AddDate(0, 0, -int(now.Weekday()+6)%7)

	// Create a slice to store the days of the week
	days := make([]string, 7)

	// Loop through the next 7 days to get the days of the week
	for i := 0; i < 7; i++ {
		day := monday.AddDate(0, 0, i)
		days[i] = day.Format("Mon")
	}

	return days
}

func StringToSate(dateString string) (time.Time, error) {

	var dueDate time.Time

	if _dueDate, err := time.Parse(time.DateOnly, dateString); err != nil {
		return dueDate, errors.New("Due date error: " + dateString)
	} else {
		dueDate = _dueDate
	}

	return dueDate, nil
}
