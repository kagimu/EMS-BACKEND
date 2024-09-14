package commons

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/services"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

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

	entityService := services.EntityService{}

	if user.Number != "" {
		currentEntity, err := entityService.GetCurrentEntity(user.Number)
		if err != nil {
			fmt.Println("Getting user current entity error: " + err.Error())
		} else {
			c.Set("entity_number", currentEntity.Number)
		}
	}

	return &user, nil
}

func GetUserEntityNumber(c *gin.Context) (string, error) {

	entityService := services.EntityService{}

	userNumber := c.Request.Header["User-Number"][0]

	fmt.Println("Getting current entity for: " + userNumber)

	var currentEntity string

	if userNumber != "" {
		checkEntity, err := entityService.GetCurrentEntity(userNumber)
		if err != nil {
			fmt.Println("Getting user current entity error: " + err.Error())
		} else {
			currentEntity = checkEntity.Number
			c.Set("entity_number", currentEntity)
		}
	}

	fmt.Println("Current entity set to: " + currentEntity)

	return currentEntity, nil
}

type PermissionGroup map[string][]string

func PermissionGroups() []PermissionGroup {

	return []PermissionGroup{
		{"CURRICULUM": []string{"create_curriculum", "view_curriculum", "export_curriculum",
			"delete_curriculum", "import_curriculum", "create_learning_area", "view_learning_area", "update_learning_area", "view_subject-cluster", "view_subject",
			"view_strands", "view_topic", "view_sub-topic", "view_concept"}},
		{"SCHEME OF WORK": []string{"view_scheme-of-work", "create_scheme", "import_scheme", "export_scheme", "delete_scheme"}},
		{"SETUP": []string{"view_school", "add_school", "edit_school", "delete_school", "create_campus", "edit_campus", "delete_campus", "create_academic-year", "view_academic-year", "update_academic-year",
			"create_term", "view_term", "update_term", "create_class", "view_class", "update_class", "create_stream", "view_streams", "update_stream", "view_campus_settings", "view_aggregate", "view_grading", "view_assessment_schedule", "view_report_card_settings", "view_student_score", "view_study_period", "view_study_period_group_names", "view_academic_year", "view_classes", "view_levels"}},
		{"USERS": []string{"create_user", "view_users", "view_roles", "update_user_details", "delete_user", "view_account-management", "view_user-management", "add_timetable", "view_timetable-details"}},
		{"PERMISSIONS": []string{"create_permission", "view_permission", "view_scopes", "delete_permission", "update_permission"}},
		{"STUDENT APPLICATION": []string{"view_applications", "view_application_list"}},
		{"STUDENTS": []string{"view_students", "update_student_details"}},
		{"CALENDAR": []string{"view_calendar", "view_calendar-details"}},
		{"INTERVIEWS": []string{"view_interview-details", "create_interview_schedules"}},
		{"INVOICING": []string{"add_custom_invoice", "pay_invoice", "cancel_invoice", "cancel_payment"}},
		{"Catering": []string{"view_catering_menu", "create_unit_of_measure", "record_meals", "view_meals", "view_unit_of_measure", "view_menu_items", "view_menu_packs", "view_feeding_groups",
			"create_feeding_groups", "view_meal_cost_projections", "view_meals_delivery", "view_weekly_meal_items",
			"view_groups_on_calendar", "view_daily_group_courses", "view_meal_course_menu_packs", "view_meal_courses",
			"view_stock", "create_new_stock", "view_stocking_history", "can_add_new_unit_of_measure", "edit_unit_of_measure", "delete_unit_of_measure", "view_ingredients", "view_menu_Items", "create_menu_item",
			"edit_menu_item", "remove_menu_item", "delete_menu_pack", "edit_menu_pack", "create_menu_pack", "delete_feeding_group", "edit_feeding_group", "view_feeding_group_list",
			"create_meal_course", "edit_meal_course", "delete_meal_course", "create_meal_course_menu_pack", "edit_meal_course_menu_pack", "delete_meal_course_menu_pack", "create_daily_group_courses", "activate_groups_on_calendar", "activate_meal_course", "serve_meal"}},
		{"PRODUCTS & SERVICES": []string{"create_product", "create_service", "edit_product", "edit_service", "view_services", "view_products"}},
		{"ASSESSMENT": []string{"view_grade-scale", "view_grade-standard", "view_assessment-list", "view_examtype", "add_assessment_schedule", "delete_assessment_schedule", "edit_assessment_schedule"}},
		{"ADMISSION": []string{"admit_student", "view_admissions", "cancel_admission"}},
		{"ENROLLMENT": []string{"enroll_student", "view_enrollments", "terminate_enrollment", "re-instate_enrollment"}},
		{"SETTINGS": []string{"access_settings_menu", "view_learning_area", "update_learning_area", "create_study_period_group_name", "view_study_period_group_name", "view_entities"}},
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
	//if extension != ".pdf" {
	//	return " ", errors.New("file should be a .pdf")
	//}

	fileName := name + extension
	filePath := "/fileRepo/" + savePath + "/" + fileName

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

func ConvertStringToDate(dateString string) (time.Time, error) {

	var dueDate time.Time

	if _dueDate, err := time.Parse(time.DateOnly, dateString); err != nil {
		return dueDate, errors.New("Due date error: " + dateString)
	} else {
		dueDate = _dueDate
	}

	return dueDate, nil
}

func GenerateDownloadPath(number string) string {
	if number == "" {
		return ""
	}
	return os.Getenv("BASE_URL") + "attachments/download/" + number
}
