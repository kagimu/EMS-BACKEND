package services

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	authResponse "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	campusRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	cateringRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/repositories"
	schoolRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

func AddMealCoursesToFeedingGroup(courseToFeedingGroupRequest requests.MealCourseToFeedingGroupRequest, currentUser string) error {

	feedingGroup, err := cateringRepositories.GetFeedingGroupByNumber(courseToFeedingGroupRequest.FeedingGroupNumber)

	if err != nil {
		return err
	}

	if feedingGroup.Number == "" {
		return errors.New("invalid feeding group number")
	}

	if feedingGroup.Status != constants.StatusActive {
		return errors.New("selected group is not active")
	}

	trn := config.DB.Begin()

	for _, day := range courseToFeedingGroupRequest.Days {

		if checkErr := commons.CheckDay(day); checkErr != "" {
			return errors.New(checkErr)
		}

		groupMealCourses, err := cateringRepositories.GetAllFeedingGroupMealCoursesForDay(feedingGroup.Number, day, trn)
		if err != nil {
			trn.Rollback()
			return errors.New("failed to get all feeding group meal courses: " + err.Error())
		}

		for _, groupMealCourse := range groupMealCourses {
			err = cateringRepositories.DeleteFeedingGroupMealCourse(groupMealCourse, trn)
			if err != nil {
				trn.Rollback()
				return errors.New("failed to clear feeding group meal course: " + err.Error())
			}
		}

		if err = saveMealCoursesForDay(courseToFeedingGroupRequest.MealCourses, day, feedingGroup, currentUser, trn); err != nil {
			trn.Rollback()
			return err
		}

	}

	if err = trn.Commit().Error; err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}

	return nil
}

func saveMealCoursesForDay(mealCourses []string, day string, feedingGroup models.FeedingGroup, currentUser string, trn *gorm.DB) error {

	for _, mealCourse := range mealCourses {
		checkMeal, err := cateringRepositories.GetMealCourseByNumber(mealCourse, trn)
		if err != nil {
			return err
		}
		if checkMeal.Number == "" {
			return errors.New(fmt.Sprintf("invalid '%s' meal ID", day))
		}

		if checkMeal.Status != constants.StatusActive {
			return errors.New(fmt.Sprintf("meal course '%s' is not active", checkMeal.Name))
		}

		groupMeal, err := cateringRepositories.GetFeedingGroupMealCourseByMealAndGroup(checkMeal.Number, day, feedingGroup.Number, trn)
		if err != nil {
			return err
		}

		if groupMeal.Number != "" {
			return errors.New(fmt.Sprintf("meal course '%s' already added to group '%s' for '%s'", checkMeal.Name, feedingGroup.Number, day))
		}

		newGroupMealCourse := models.FeedingGroupMealCourse{
			Number:             uuid.NewV4().String(),
			FeedingGroupNumber: feedingGroup.Number,
			MealCourseNumber:   mealCourse,
			Day:                day,
			UserNumber:         currentUser,
		}

		if err := cateringRepositories.SaveFeedingGroupMealCourse(newGroupMealCourse, trn); err != nil {
			return err
		}

	}

	return nil
}

func GetFeedingGroupMealCourseCalendar(schoolNumber string) (responses.FeedingGroupMealCourseCalendar, error) {

	var response responses.FeedingGroupMealCourseCalendar

	school, err := schoolRepositories.GetSchoolByNumber(schoolNumber)

	if err != nil {
		return response, err
	}

	if school.Number == "" {
		return response, errors.New("invalid school number")
	}

	headers := []string{" "}

	for _, day := range commons.DaysOfTheWeek() {
		headers = append(headers, day)
	}

	groups, err := cateringRepositories.GetAllFeedingGroups(school.Number, constants.StatusActive)

	if err != nil {
		return response, err
	}

	var items []responses.GroupCalendarRawResponse

	for _, group := range groups {

		dailyMenu := getGroupDailyMenu(group.Number)

		item := responses.GroupCalendarRawResponse{
			GroupName:   group.Name,
			MealCourses: dailyMenu,
		}

		items = append(items, item)

	}

	response.FeedingGroup = items

	response.Headers = headers

	return response, nil
}

func getGroupDailyMenu(groupNumber string) []responses.GroupDayMealCourses {

	var weeklyItems []responses.GroupDayMealCourses

	for _, day := range commons.DaysOfTheWeek() {

		var item responses.GroupDayMealCourses
		item.Day = day
		item.Active = false

		mealCourses, err := cateringRepositories.GetAllFeedingGroupMealCoursesForDay(groupNumber, day)

		var meals []string

		if err != nil {
			item.Error = err.Error()
		} else {

			for _, _mealCourse := range mealCourses {
				mealCourse := _mealCourse.MealCourse
				if mealCourse.Status == constants.StatusActive {
					meals = append(meals, mealCourse.Name)
				}
			}

			if len(mealCourses) > 0 && len(meals) == 0 {
				item.Error = "Inactive meal courses"
			} else {
				item.Error = "No meal course(s) set"
			}

		}

		item.MealCourses = meals

		if len(meals) > 0 {
			item.Active = true
		}

		weeklyItems = append(weeklyItems, item)

	}

	return weeklyItems
}

func GetUserMeals(mealCourseNumber string, campusNumber string, feedingGroupNumber string, weekDay ...string) ([]responses.UserMealResponse, error) {

	mealsResponse := make([]responses.UserMealResponse, 0)

	currentTime := time.Now()
	day := currentTime.Format("Mon")

	if len(weekDay) > 0 {
		day = weekDay[0]
	} else {
		checkDay := commons.CheckDay(day)
		if checkDay != "" {
			return mealsResponse, errors.New(checkDay)
		}
	}

	campus, err := campusRepositories.GetCampusByNumber(campusNumber)

	if err != nil {
		return mealsResponse, err
	}

	if campus.Number == "" {
		return mealsResponse, errors.New("invalid campus number")
	}

	mealCourse, err := cateringRepositories.GetMealCourseByNumber(mealCourseNumber)
	if err != nil {
		return mealsResponse, err
	}

	if mealCourse.Status != constants.StatusActive {
		return mealsResponse, errors.New("meal course is: '" + mealCourse.Status + "'")
	}

	checkSchool, err := cateringRepositories.GetMealCourseByNumberAndSchool(mealCourse.Number, campus.SchoolNumber)
	if err != nil {
		return mealsResponse, err
	}

	if checkSchool.Number == "" {
		return mealsResponse, errors.New("meal course can not be used on the selected campus")
	}

	mealCourseFeedingSession, err := cateringRepositories.GetFeedingSessionForCampusAndMealCourseAndDay(day, mealCourseNumber, campusNumber)
	if err != nil {
		return mealsResponse, err
	}

	if mealCourseFeedingSession.Number == "" {
		return mealsResponse, errors.New("menu not set up for this meal course on '" + day + "'")
	}

	feedingSession := mealCourseFeedingSession.FeedingSession

	feedItems, err := cateringRepositories.GetFeedingSessionItems(feedingSession.Number)
	if err != nil {
		return mealsResponse, err
	}

	if len(feedItems) == 0 {
		return mealsResponse, errors.New("no items set to '" + feedingSession.Number + "' session")
	}

	var feedingGroups []models.FeedingGroup

	if feedingGroupNumber != "" {

		checkGroup, err := cateringRepositories.GetFeedingGroupByNumber(feedingGroupNumber)
		if err != nil {
			return mealsResponse, err
		}

		if checkGroup.Number == "" || checkGroup.SchoolNumber != campus.SchoolNumber {
			return mealsResponse, errors.New("invalid group number")
		}

		checkAllowedToHaveMeal, err := cateringRepositories.GetFeedingGroupMealCourseByMealAndGroup(mealCourseNumber, day, feedingGroupNumber)
		if err != nil {
			return mealsResponse, err
		}

		if checkAllowedToHaveMeal.Number == "" {
			return mealsResponse, errors.New("this group should not be served today")
		}

		feedingGroups = append(feedingGroups, checkGroup)

	} else {

		allFeedingGroupsForMealAndDay, err := cateringRepositories.GetAllFeedingGroupMealCourseByMealAndDay(mealCourseNumber, day)
		if err != nil {
			return mealsResponse, err
		}

		for _, feedingGroup := range allFeedingGroupsForMealAndDay {
			feedingGroups = append(feedingGroups, feedingGroup.FeedingGroup)
		}

	}

	for _, feedingGroup := range feedingGroups {

		allUsers, err := cateringRepositories.GetFeedingGroupUserByGroupNumber("", "", feedingGroup.Number)
		if err != nil {
			return mealsResponse, err
		}

		for _, _user := range allUsers {

			user := _user.User

			userMeal := responses.UserMealResponse{
				UserNumber:           user.Number,
				Name:                 user.FirstName,
				FeedingGroupNumber:   feedingGroup.Number,
				FeedingGroupName:     feedingGroup.Name,
				FeedingSessionNumber: feedingSession.Number,
			}

			meals := findUserMeal(user.Number, feedItems, mealCourse.Number)

			for _, meal := range meals {
				if meal.ServedQuantity != 0 {
					userMeal.Served = true
				}
			}

			userMeal.Meals = meals

			mealsResponse = append(mealsResponse, userMeal)
		}

	}

	return mealsResponse, nil
}

// GetLearnerMeals is a variant of GetUserMeals that restricts results to enrolled learners only
// and enriches each entry with full enrollment details (class, level, study period, academic year, etc.).
// Supported filters: mealCourseNumber (required), campusNumber (required), feedingGroupNumber (optional),
// studyPeriodNumber (optional), academicYearNumber (optional), campusClassNumber (optional), streamNumber (optional), classNumber (optional).
func GetLearnerMeals(mealCourseNumber, campusNumber, feedingGroupNumber, studyPeriodNumber, academicYearNumber, campusClassNumber, streamNumber string, classNumber string, weekDay ...string) ([]responses.LearnerMealResponse, error) {

	mealsResponse := make([]responses.LearnerMealResponse, 0)

	currentTime := time.Now()
	day := currentTime.Format("Mon")

	if len(weekDay) > 0 {
		day = weekDay[0]
	} else {
		checkDay := commons.CheckDay(day)
		if checkDay != "" {
			return mealsResponse, errors.New(checkDay)
		}
	}

	campus, err := campusRepositories.GetCampusByNumber(campusNumber)
	if err != nil {
		return mealsResponse, err
	}
	if campus.Number == "" {
		return mealsResponse, errors.New("invalid campus number")
	}

	mealCourse, err := cateringRepositories.GetMealCourseByNumber(mealCourseNumber)
	if err != nil {
		return mealsResponse, err
	}
	if mealCourse.Status != constants.StatusActive {
		return mealsResponse, errors.New("meal course is: '" + mealCourse.Status + "'")
	}

	checkSchool, err := cateringRepositories.GetMealCourseByNumberAndSchool(mealCourse.Number, campus.SchoolNumber)
	if err != nil {
		return mealsResponse, err
	}
	if checkSchool.Number == "" {
		return mealsResponse, errors.New("meal course can not be used on the selected campus")
	}

	mealCourseFeedingSession, err := cateringRepositories.GetFeedingSessionForCampusAndMealCourseAndDay(day, mealCourseNumber, campusNumber)
	if err != nil {
		return mealsResponse, err
	}
	if mealCourseFeedingSession.Number == "" {
		return mealsResponse, errors.New("menu not set up for this meal course on '" + day + "'")
	}

	feedingSession := mealCourseFeedingSession.FeedingSession

	feedItems, err := cateringRepositories.GetFeedingSessionItems(feedingSession.Number)
	if err != nil {
		return mealsResponse, err
	}
	if len(feedItems) == 0 {
		return mealsResponse, errors.New("no items set to '" + feedingSession.Number + "' session")
	}

	var feedingGroups []models.FeedingGroup

	if feedingGroupNumber != "" {
		checkGroup, err := cateringRepositories.GetFeedingGroupByNumber(feedingGroupNumber)
		if err != nil {
			return mealsResponse, err
		}
		if checkGroup.Number == "" || checkGroup.SchoolNumber != campus.SchoolNumber {
			return mealsResponse, errors.New("invalid group number")
		}
		checkAllowedToHaveMeal, err := cateringRepositories.GetFeedingGroupMealCourseByMealAndGroup(mealCourseNumber, day, feedingGroupNumber)
		if err != nil {
			return mealsResponse, err
		}
		if checkAllowedToHaveMeal.Number == "" {
			return mealsResponse, errors.New("this group should not be served today")
		}
		feedingGroups = append(feedingGroups, checkGroup)
	} else {
		allFeedingGroupsForMealAndDay, err := cateringRepositories.GetAllFeedingGroupMealCourseByMealAndDay(mealCourseNumber, day)
		if err != nil {
			return mealsResponse, err
		}
		for _, fg := range allFeedingGroupsForMealAndDay {
			feedingGroups = append(feedingGroups, fg.FeedingGroup)
		}
	}

	db := config.DB

	for _, feedingGroup := range feedingGroups {

		// Query learners that are in this feeding group AND have an active enrollment,
		// joined with all enrollment/class/study-period details.
		var learnerRows []struct {
			UserNumber           string
			Gender               string
			FirstName            string
			LastName             string
			Surname              string
			SecondName           string
			OtherName            string
			Photo                string
			Adn                  string
			EnrollmentNumber     string
			EnrollmentStatus     string
			EnrollmentDate       string
			ContinuingStudent    bool
			CampusClassNumber    string
			StudyPeriodNumber    string
			StreamNumber         string
			Class                string
			Level                string
			AcademicYear         string
			StudyPeriod          string
			StudyPeriodStart     string
			StudyPeriodEnd       string
		}

		q := db.Table("feeding_group_users").
			Select("feeding_group_users.user_number, users.gender, users.first_name, users.last_name, "+
				"students.surname, students.second_name, students.other_name, students.photo, students.adn, "+
				"enrollments.number as enrollment_number, enrollments.status as enrollment_status, "+
				"enrollments.enrollment_date as enrollment_date, enrollments.continuing_student, "+
				"enrollments.campus_class_number, enrollments.study_period_number, enrollments.stream_number, "+
				"classes.name as class, levels.name as level, "+
				"academic_years.academic_year_name as academic_year, study_period_names.name as study_period, "+
				"study_periods.study_period_start_date as study_period_start, "+
				"study_periods.study_period_end_date as study_period_end").
			Joins("JOIN users ON feeding_group_users.user_number = users.number").
			Joins("JOIN students ON students.user_number = users.number").
			Joins("JOIN enrollments ON enrollments.student_number = students.number AND enrollments.status = 'ACTIVE'").
			Joins("LEFT JOIN campus_classes ON campus_classes.number = enrollments.campus_class_number").
			Joins("LEFT JOIN curriculum_classes ON curriculum_classes.number = campus_classes.curriculum_class_number").
			Joins("LEFT JOIN classes ON classes.number = curriculum_classes.class_number").
			Joins("LEFT JOIN levels ON levels.number = classes.level_number").
			Joins("LEFT JOIN study_periods ON study_periods.number = enrollments.study_period_number").
			Joins("LEFT JOIN study_period_names ON study_period_names.number = study_periods.study_period_name_number").
			Joins("LEFT JOIN academic_years ON academic_years.number = study_periods.academic_year_number").
			Where("feeding_group_users.feeding_group_number = ?", feedingGroup.Number).
			Where("feeding_group_users.campus_number = ?", campusNumber).
			Where("users.type = ?", "Student")

		if studyPeriodNumber != "" {
			q = q.Where("enrollments.study_period_number = ?", studyPeriodNumber)
		}
		if academicYearNumber != "" {
			q = q.Where("study_periods.academic_year_number = ?", academicYearNumber)
		}
		if campusClassNumber != "" {
			q = q.Where("enrollments.campus_class_number = ?", campusClassNumber)
		}
		if streamNumber != "" {
			q = q.Where("enrollments.stream_number = ?", streamNumber)
		}
		if classNumber != "" {
			q = q.Where("classes.number = ?", classNumber)
		}

		if err := q.Scan(&learnerRows).Error; err != nil {
			return mealsResponse, err
		}

		for _, row := range learnerRows {

			learnerMeal := responses.LearnerMealResponse{
				// Meal / session fields
				FeedingGroupNumber:   feedingGroup.Number,
				FeedingGroupName:     feedingGroup.Name,
				FeedingSessionNumber: feedingSession.Number,
				// Learner identity
				UserNumber:        row.UserNumber,
				Name:              row.FirstName + " " + row.LastName,
				StudentSurname:    row.Surname,
				StudentSecondName: row.SecondName,
				StudentOtherName:  row.OtherName,
				Gender:            row.Gender,
				ADN:               row.Adn,
				Photo:             commons.GenerateDownloadPath(row.Photo, true),
				// Enrollment details
				EnrollmentNumber:     row.EnrollmentNumber,
				EnrollmentStatus:     row.EnrollmentStatus,
				EnrollmentDate:       row.EnrollmentDate,
				ContinuingStudent:    row.ContinuingStudent,
				StudyPeriodStartDate: row.StudyPeriodStart,
				StudyPeriodEndDate:   row.StudyPeriodEnd,
				CampusClassNumber:    row.CampusClassNumber,
				StudyPeriodNumber:    row.StudyPeriodNumber,
				StreamNumber:         row.StreamNumber,
				Class:                row.Class,
				Level:                row.Level,
				AcademicYear:         row.AcademicYear,
				StudyPeriod:          row.StudyPeriod,
				// Campus / school
				CampusNumber: campusNumber,
				CampusName:   campus.Name,
				SchoolNumber: campus.SchoolNumber,
			}

			meals := findUserMeal(row.UserNumber, feedItems, mealCourse.Number)

			for _, meal := range meals {
				if meal.ServedQuantity != 0 {
					learnerMeal.Served = true
				}
			}
			learnerMeal.Meals = meals

			mealsResponse = append(mealsResponse, learnerMeal)
		}
	}

	return mealsResponse, nil
}

// GetLearnerMealsHistory retrieves the meal delivery history for enrolled learners on a specific date.
func GetLearnerMealsHistory(mealCourseNumber, campusNumber, feedingGroupNumber, studyPeriodNumber, academicYearNumber, campusClassNumber, streamNumber, classNumber, queryDate string) ([]responses.LearnerMealHistoryResponse, error) {

	historyResponse := make([]responses.LearnerMealHistoryResponse, 0)

	targetDate, err := time.Parse(time.DateOnly, queryDate)
	if err != nil {
		return historyResponse, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	day := targetDate.Format("Mon")
	checkDay := commons.CheckDay(day)
	if checkDay != "" {
		return historyResponse, errors.New("selected date is not a valid feeding day: " + checkDay)
	}

	campus, err := campusRepositories.GetCampusByNumber(campusNumber)
	if err != nil {
		return historyResponse, err
	}
	if campus.Number == "" {
		return historyResponse, errors.New("invalid campus number")
	}

	mealCourse, err := cateringRepositories.GetMealCourseByNumber(mealCourseNumber)
	if err != nil {
		return historyResponse, err
	}
	if mealCourse.Status != constants.StatusActive {
		return historyResponse, errors.New("meal course is: '" + mealCourse.Status + "'")
	}

	checkSchool, err := cateringRepositories.GetMealCourseByNumberAndSchool(mealCourse.Number, campus.SchoolNumber)
	if err != nil {
		return historyResponse, err
	}
	if checkSchool.Number == "" {
		return historyResponse, errors.New("meal course can not be used on the selected campus")
	}

	mealCourseFeedingSession, err := cateringRepositories.GetFeedingSessionForCampusAndMealCourseAndDay(day, mealCourseNumber, campusNumber)
	if err != nil {
		return historyResponse, err
	}
	if mealCourseFeedingSession.Number == "" {
		return historyResponse, errors.New("menu not set up for this meal course on '" + day + "'")
	}

	var feedingGroups []models.FeedingGroup

	if feedingGroupNumber != "" {
		checkGroup, err := cateringRepositories.GetFeedingGroupByNumber(feedingGroupNumber)
		if err != nil {
			return historyResponse, err
		}
		if checkGroup.Number == "" || checkGroup.SchoolNumber != campus.SchoolNumber {
			return historyResponse, errors.New("invalid group number")
		}
		checkAllowedToHaveMeal, err := cateringRepositories.GetFeedingGroupMealCourseByMealAndGroup(mealCourseNumber, day, feedingGroupNumber)
		if err != nil {
			return historyResponse, err
		}
		if checkAllowedToHaveMeal.Number == "" {
			return historyResponse, errors.New("this group should not be served on this date")
		}
		feedingGroups = append(feedingGroups, checkGroup)
	} else {
		allFeedingGroupsForMealAndDay, err := cateringRepositories.GetAllFeedingGroupMealCourseByMealAndDay(mealCourseNumber, day)
		if err != nil {
			return historyResponse, err
		}
		for _, fg := range allFeedingGroupsForMealAndDay {
			feedingGroups = append(feedingGroups, fg.FeedingGroup)
		}
	}

	db := config.DB

	for _, feedingGroup := range feedingGroups {

		var learnerRows []struct {
			UserNumber string
			FirstName  string
			LastName   string
		}

		q := db.Table("feeding_group_users").
			Select("feeding_group_users.user_number, users.first_name, users.last_name").
			Joins("JOIN users ON feeding_group_users.user_number = users.number").
			Joins("JOIN students ON students.user_number = users.number").
			Joins("JOIN enrollments ON enrollments.student_number = students.number AND enrollments.status = 'ACTIVE'").
			Joins("LEFT JOIN campus_classes ON campus_classes.number = enrollments.campus_class_number").
			Joins("LEFT JOIN curriculum_classes ON curriculum_classes.number = campus_classes.curriculum_class_number").
			Joins("LEFT JOIN classes ON classes.number = curriculum_classes.class_number").
			Joins("LEFT JOIN study_periods ON study_periods.number = enrollments.study_period_number").
			Where("feeding_group_users.feeding_group_number = ?", feedingGroup.Number).
			Where("feeding_group_users.campus_number = ?", campusNumber).
			Where("users.type = ?", "Student")

		if studyPeriodNumber != "" {
			q = q.Where("enrollments.study_period_number = ?", studyPeriodNumber)
		}
		if academicYearNumber != "" {
			q = q.Where("study_periods.academic_year_number = ?", academicYearNumber)
		}
		if campusClassNumber != "" {
			q = q.Where("enrollments.campus_class_number = ?", campusClassNumber)
		}
		if streamNumber != "" {
			q = q.Where("enrollments.stream_number = ?", streamNumber)
		}
		if classNumber != "" {
			q = q.Where("classes.number = ?", classNumber)
		}

		if err := q.Scan(&learnerRows).Error; err != nil {
			return historyResponse, err
		}

		for _, row := range learnerRows {
			
			mealRecord, err := cateringRepositories.GetMealByMealNumberAndDateAndUserNumber(mealCourse.Number, targetDate, row.UserNumber)
			if err != nil {
				return historyResponse, err
			}

			historyItem := responses.LearnerMealHistoryResponse{
				Name:             row.FirstName + " " + row.LastName,
				FeedingGroupName: feedingGroup.Name,
				MealName:         mealCourse.Name,
			}

			if mealRecord.Number != "" {
				historyItem.Served = true
				historyItem.Date = mealRecord.Date.Format(time.DateOnly)
				historyItem.Time = mealRecord.Time
			} else {
				historyItem.Served = false
				historyItem.Date = queryDate
				historyItem.Time = "-"
			}

			historyResponse = append(historyResponse, historyItem)
		}
	}

	return historyResponse, nil
}

// GetStaffMeals is the staff counterpart of GetLearnerMeals.
// It returns staff members that belong to feeding groups for the given school and meal course,
// together with their served status for today's meal.
//
// Required: mealCourseNumber, schoolNumber
// Optional: feedingGroupNumber
func GetStaffMeals(mealCourseNumber, schoolNumber, feedingGroupNumber string) ([]responses.StaffMealResponse, error) {

	mealsResponse := make([]responses.StaffMealResponse, 0)

	// --- 1. Resolve current day-of-week ---
	currentTime := time.Now()
	day := currentTime.Format("Mon")
	if checkDay := commons.CheckDay(day); checkDay != "" {
		return mealsResponse, errors.New(checkDay)
	}

	// --- 2. Validate school ---
	school, err := schoolRepositories.GetSchoolByNumber(schoolNumber)
	if err != nil {
		return mealsResponse, err
	}
	if school.Number == "" {
		return mealsResponse, errors.New("invalid school number")
	}

	// --- 3. Validate meal course (must be active and belong to this school) ---
	mealCourse, err := cateringRepositories.GetMealCourseByNumber(mealCourseNumber)
	if err != nil {
		return mealsResponse, err
	}
	if mealCourse.Status != constants.StatusActive {
		return mealsResponse, errors.New("meal course is: '" + mealCourse.Status + "'")
	}
	checkSchool, err := cateringRepositories.GetMealCourseByNumberAndSchool(mealCourse.Number, schoolNumber)
	if err != nil {
		return mealsResponse, err
	}
	if checkSchool.Number == "" {
		return mealsResponse, errors.New("meal course does not belong to the selected school")
	}

	// --- 4. Resolve feeding groups ---
	var feedingGroups []models.FeedingGroup

	if feedingGroupNumber != "" {
		checkGroup, err := cateringRepositories.GetFeedingGroupByNumber(feedingGroupNumber)
		if err != nil {
			return mealsResponse, err
		}
		if checkGroup.Number == "" || checkGroup.SchoolNumber != schoolNumber {
			return mealsResponse, errors.New("invalid group number")
		}
		checkAllowedToHaveMeal, err := cateringRepositories.GetFeedingGroupMealCourseByMealAndGroup(mealCourseNumber, day, feedingGroupNumber)
		if err != nil {
			return mealsResponse, err
		}
		if checkAllowedToHaveMeal.Number == "" {
			return mealsResponse, errors.New("this group should not be served today")
		}
		feedingGroups = append(feedingGroups, checkGroup)
	} else {
		allFeedingGroupsForMealAndDay, err := cateringRepositories.GetAllFeedingGroupMealCourseByMealAndDay(mealCourseNumber, day)
		if err != nil {
			return mealsResponse, err
		}
		fmt.Printf("[GetStaffMeals] day=%s, meal=%s → %d feeding group meal courses found\n", day, mealCourseNumber, len(allFeedingGroupsForMealAndDay))
		for _, fg := range allFeedingGroupsForMealAndDay {
			fmt.Printf("[GetStaffMeals]   group=%s schoolNumber=%s (want %s)\n", fg.FeedingGroup.Number, fg.FeedingGroup.SchoolNumber, schoolNumber)
			if fg.FeedingGroup.SchoolNumber == schoolNumber {
				feedingGroups = append(feedingGroups, fg.FeedingGroup)
			}
		}
		fmt.Printf("[GetStaffMeals] feeding groups after school filter: %d\n", len(feedingGroups))
	}

	// --- 5. Process each feeding group ---
	db := config.DB

	for _, feedingGroup := range feedingGroups {

		// 5a. Fetch all staff assigned to this group with full deployment preloads.
		var staffRows []models.FeedingGroupStaff
		if err := db.
			Preload("Campus").
			Preload("FeedingGroup").
			Preload("StaffDeployment").
			Preload("StaffDeployment.JobOffer").
			Preload("StaffDeployment.JobOffer.JobApplication").
			Preload("StaffDeployment.JobOffer.JobInterview").
			Preload("StaffDeployment.JobOffer.JobInterview.JobVacancy").
			Preload("StaffDeployment.JobOffer.JobInterview.JobVacancy.JobPosition").
			Preload("StaffDeployment.JobOffer.JobInterview.JobVacancy.Campus").
			Where("feeding_group_number = ?", feedingGroup.Number).
			Find(&staffRows).Error; err != nil {
			return mealsResponse, err
		}

		fmt.Printf("[GetStaffMeals] group=%s → %d staff rows\n", feedingGroup.Number, len(staffRows))
		if len(staffRows) == 0 {
			continue
		}

		// 5b. Resolve the feeding session using the campus stored on the first staff record
		//     for this group (all staff in a group share the same campus).
		campusNumber := staffRows[0].CampusNumber
		campus := staffRows[0].Campus

		fmt.Printf("[GetStaffMeals] group=%s → campus=%s, looking for session day=%s meal=%s\n", feedingGroup.Number, campusNumber, day, mealCourseNumber)

		mealCourseFeedingSession, err := cateringRepositories.GetFeedingSessionForCampusAndMealCourseAndDay(day, mealCourseNumber, campusNumber)
		if err != nil {
			return mealsResponse, err
		}
		if mealCourseFeedingSession.Number == "" {
			fmt.Printf("[GetStaffMeals] group=%s → NO feeding session found for campus=%s day=%s — skipping\n", feedingGroup.Number, campusNumber, day)
			continue
		}

		feedingSession := mealCourseFeedingSession.FeedingSession

		feedItems, err := cateringRepositories.GetFeedingSessionItems(feedingSession.Number)
		if err != nil {
			return mealsResponse, err
		}

		// 5c. Build a response record per staff member.
		for _, s := range staffRows {
			dep := s.StaffDeployment
			app := dep.JobOffer.JobApplication

			// Resolve display name.
			name := app.FirstName
			if app.LastName != "" {
				name += " " + app.LastName
			}
			if name == "" {
				name = "Staff Member"
			}

			// Resolve job position.
			position := ""
			if dep.JobOffer.JobInterview.JobVacancy.JobPosition.PostName != "" {
				position = dep.JobOffer.JobInterview.JobVacancy.JobPosition.PostName
			}

			// Resolve campus name (fall back to vacancy campus if own campus is missing).
			campusName := campus.Name
			if campusName == "" {
				campusName = dep.JobOffer.JobInterview.JobVacancy.Campus.Name
			}

			// Resolve user_number for served-status look-up.
			// Staff deployments link to a user account; fall back to application's user.
			userNumber := dep.UserNumber
			if userNumber == "" {
				userNumber = app.UserNumber
			}

			staffMeal := responses.StaffMealResponse{
				UserNumber:           userNumber,
				Name:                 name,
				FeedingGroupNumber:   feedingGroup.Number,
				FeedingGroupName:     feedingGroup.Name,
				FeedingSessionNumber: feedingSession.Number,
				Gender:               "",
				STN:                  dep.StaffIdNumber,
				Email:                app.Email,
				PhoneNumber:          app.PhoneNumber,
				JobPosition:          position,
				SchoolNumber:         schoolNumber,
				CampusNumber:         campusNumber,
				CampusName:           campusName,
			}

			// Check served status (only meaningful when user_number is linked).
			if userNumber != "" {
				meals := findUserMeal(userNumber, feedItems, mealCourse.Number)
				for _, meal := range meals {
					if meal.ServedQuantity != 0 {
						staffMeal.Served = true
					}
				}
				staffMeal.Meals = meals
			}

			mealsResponse = append(mealsResponse, staffMeal)
		}
	}

	return mealsResponse, nil
}

func findUserMeal(userNumber string, feedingSessionItems []models.FeedingSessionItem, mealCourse string) []responses.SingleUserMealItem {

	var items []responses.SingleUserMealItem

	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")

	date_, err := time.Parse(time.DateOnly, currentDate)
	if err != nil {
		fmt.Println("Error formatting today's date: " + currentDate)
	}

	for _, feedingSessionItem := range feedingSessionItems {

		menuItem, err := cateringRepositories.GetMenuItemByNumber(feedingSessionItem.MenuItemNumber)
		if err != nil {
			fmt.Println("Error getting menu item from feeding session")
			return items
		}

		unitOfMeasure, err := cateringRepositories.GetUnitOfMeasureByNumber(menuItem.UnitOfMeasureNumber)
		if err != nil {
			fmt.Println("Error getting unit of measure" + err.Error())
			unitOfMeasure = models.UnitOfMeasure{}
		}

		unitOfMeasureName := "None"
		if unitOfMeasure.Number != "" {
			unitOfMeasureName = unitOfMeasure.Name
		}

		var served float64

		servedStock, err := cateringRepositories.GetMealStockingByMealNumberAndDateUserNumber(mealCourse, date_, userNumber, menuItem.Number)

		if err != nil {
			fmt.Println("Error getting menu item served quantity: " + err.Error())
			return items
		}

		if servedStock.Number != "" {
			served = -1 * servedStock.Quantity
		} else {
			served = 0
		}

		singleUserMenu := responses.SingleUserMealItem{
			Quantity:       feedingSessionItem.Quantity,
			ServedQuantity: served,
			MenuItemNumber: menuItem.Number,
			MenuItemName:   menuItem.Name,
			UnitOfMeasure:  unitOfMeasureName,
		}

		items = append(items, singleUserMenu)
	}

	return items
}

func GetFeedingGroupsBySchoolNumber(schoolNumber string, status string) ([]responses.FeedingGroupResponse, error) {

	var response []responses.FeedingGroupResponse

	feedingGroups, err := cateringRepositories.GetAllFeedingGroups(schoolNumber, status)
	if err != nil {
		return response, err
	}

	for _, feedingGroup := range feedingGroups {

		user, err := repository.GetUserByNumber(feedingGroup.UserNumber)
		if err != nil {
			return response, err
		}

		if user.Number == "" {
			return response, errors.New("user not found")
		}

		school, err := schoolRepositories.GetSchoolByNumber(feedingGroup.SchoolNumber)
		if err != nil {
			return response, err
		}

		if school.Number == "" {
			return response, errors.New("school not found")
		}

		response = append(response, responses.FeedingGroupResponse{
			Number:       feedingGroup.Number,
			Name:         feedingGroup.Name,
			UserName:     user.FirstName + " " + user.LastName,
			Status:       feedingGroup.Status,
			SchoolNumber: school.Number,
			SchoolName:   school.Name,
		})

	}

	return response, nil
}

func AddMeal(mealRequest requests.MealRequest, currentUser string) error {

	formatDate, err := time.Parse(time.DateOnly, mealRequest.Date)
	if err != nil {
		return err
	}

	day := formatDate.Format("Mon")
	checkDay := commons.CheckDay(day)
	if checkDay != "" {
		return errors.New(checkDay)
	}

	newMeal, err := cateringRepositories.GetMealByMealNumberAndDateAndUserNumber(mealRequest.MealCourseNumber, formatDate, mealRequest.UserNumber)
	if err != nil {
		return errors.New("Checking if item served error: " + err.Error())
	}

	groupUser, err := cateringRepositories.GetFeedingGroupUserByUserNumber(mealRequest.UserNumber)
	if err != nil {
		return err
	}

	if groupUser.Number == "" {
		return errors.New("user not added to any group")
	}

	sessionItems, err := cateringRepositories.GetFeedingSessionItems(mealRequest.FeedingSessionNumber)
	if err != nil {
		return errors.New("Fetching session items error: " + err.Error())
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		for _, sessionItem := range sessionItems {
			err = processSingleItem(tx, day, formatDate, sessionItem, mealRequest, newMeal, groupUser, currentUser)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func processSingleItem(tx *gorm.DB, day string, formatDate time.Time, feedingSessionItem models.FeedingSessionItem, mealRequest requests.MealRequest, newMeal models.Meal, groupUser models.FeedingGroupUser, currentUser string) error {

	currentStock, err := cateringRepositories.GetMenuItemCurrentStock(feedingSessionItem.MenuItemNumber, tx)
	if err != nil {
		return err
	}

	if currentStock < feedingSessionItem.Quantity {
		menu, _ := cateringRepositories.GetMenuItemByNumber(feedingSessionItem.MenuItemNumber, tx)
		return errors.New("not enough stock for '" + menu.Name + "'")
	}

	if newMeal.Number == "" {
		newMeal.Number = uuid.NewV4().String()
		newMeal.MealCourseNumber = mealRequest.MealCourseNumber
		newMeal.FeedingSessionNumber = feedingSessionItem.FeedingSessionNumber
		newMeal.UserNumber = mealRequest.UserNumber
		newMeal.Date = formatDate
		newMeal.Time = mealRequest.Time
		newMeal.AddedByNumber = currentUser
		newMeal.StaffNumber = currentUser
		newMeal.Day = day
		newMeal.CampusNumber = groupUser.CampusNumber
		newMeal.FeedingGroupNumber = groupUser.FeedingGroupNumber
	} else {
		newMeal.StaffNumber = currentUser
		newMeal.CampusNumber = groupUser.CampusNumber
		newMeal.FeedingSessionNumber = feedingSessionItem.FeedingSessionNumber
		newMeal.FeedingGroupNumber = groupUser.FeedingGroupNumber
	}

	err = cateringRepositories.SaveMeal(newMeal, tx)
	if err != nil {
		return errors.New("error creating meal: " + err.Error())
	}

	stocking, err := cateringRepositories.GetMealStockByMealNumberAndMenuItemNumber(newMeal.Number, feedingSessionItem.MenuItemNumber, tx)
	if err != nil {
		return err
	}

	if stocking.Number == "" {
		stocking.Number = uuid.NewV4().String()
		stocking.MenuItemNumber = feedingSessionItem.MenuItemNumber
	}

	item, err := cateringRepositories.GetMenuItemByNumber(feedingSessionItem.MenuItemNumber, tx)
	if err != nil {
		return errors.New("Error checking item " + err.Error())
	}

	stocking.Quantity = feedingSessionItem.Quantity * -1
	stocking.CampusNumber = groupUser.CampusNumber
	stocking.UserNumber = currentUser
	stocking.PurchaseUnitOfMeasureNumber = item.PurchaseUnitOfMeasureNumber
	stocking.UnitOfMeasureNumber = item.UnitOfMeasureNumber

	if item.PurchaseUnitOfMeasureNumber == "" {
		stocking.PurchaseUnitOfMeasureNumber = item.UnitOfMeasureNumber
	}

	if stocking.PurchaseUnitOfMeasureNumber == "" {
		return errors.New("Set a base unit of measure for " + item.Name)
	}

	err = cateringRepositories.SaveStocking(stocking, tx)
	if err != nil {
		return errors.New("error reducing stock: " + err.Error())
	}

	mealStocking, err := cateringRepositories.GetMealStockingByMealNumberAndStockingNumber(newMeal.Number, stocking.Number, tx)
	if err != nil {
		return err
	}

	if mealStocking.MealNumber == "" {
		mealStocking.MealNumber = newMeal.Number
		mealStocking.StockingNumber = stocking.Number
	}

	err = cateringRepositories.SaveMealStocking(mealStocking, tx)
	if err != nil {
		return errors.New("error mapping stock with meal taken: " + err.Error())
	}

	return nil
}

func UpdateFeedingGroup(payload requests.UpdateFeedingGroupRequest, currentUserNumber string) error {

	feedingGroup, err := cateringRepositories.GetFeedingGroupByNumber(payload.FeedingGroupNumber)
	if err != nil {
		return err
	}

	if feedingGroup.Number == "" {
		return errors.New("invalid feeding group number")
	}

	checkFeedingGroupDup, err := cateringRepositories.GetFeedingGroupByName(payload.Name, feedingGroup.SchoolNumber)
	if err != nil {
		return err
	}

	if checkFeedingGroupDup.Number != "" && checkFeedingGroupDup.Number != feedingGroup.Number {
		return errors.New("this name is already defined")
	}

	feedingGroup.Status = payload.Status
	feedingGroup.Name = payload.Name
	feedingGroup.UserNumber = currentUserNumber

	err = cateringRepositories.SaveFeedingGroup(feedingGroup)
	if err != nil {
		return err
	}

	return nil
}

func CreateFeedingGroup(payload requests.FeedingGroupRequest, currentUserNumber string) error {

	checkName, err := cateringRepositories.GetFeedingGroupByName(payload.Name, payload.SchoolNumber)
	if err != nil {
		return err
	}
	if checkName.Number != "" {
		return errors.New("group with same name already defined")
	}

	newFeedingGroup := models.FeedingGroup{
		Number:       uuid.NewV4().String(),
		Name:         payload.Name,
		SchoolNumber: payload.SchoolNumber,
		Status:       payload.Status,
		UserNumber:   currentUserNumber,
	}

	err = cateringRepositories.SaveFeedingGroup(newFeedingGroup)
	if err != nil {
		return err
	}

	return nil
}

func SaveFeedingSession(payload requests.FeedingSessionRequest, currentUserNumber string) error {

	checkFeedingSession, err := cateringRepositories.GetFeedingSessionByNameAndSchool(payload.Name, payload.SchoolNumber)
	if err != nil {
		return err
	}

	if checkFeedingSession.Number != "" {
		return errors.New("this name is already defined")
	}

	err = config.DB.Transaction(func(trn *gorm.DB) error {

		newFeedingSession := models.FeedingSession{
			Number:       uuid.NewV4().String(),
			Name:         payload.Name,
			Status:       payload.Status,
			SchoolNumber: payload.SchoolNumber,
			UserNumber:   currentUserNumber,
		}

		err = cateringRepositories.SaveFeedingSession(newFeedingSession, trn)
		if err != nil {
			return err
		}

		err = processMenuItems(newFeedingSession, payload.MenuItems, currentUserNumber, trn)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func processMenuItems(menuPack models.FeedingSession, menuItemsAndQuantities []requests.MenuItemQuantity, currentUser string, tx *gorm.DB) error {
	err := checkDuplicates(menuItemsAndQuantities)
	if err != nil {
		return err
	}

	for _, menuItemQuantity := range menuItemsAndQuantities {

		saveItemAssembly := models.FeedingSessionItem{
			Number:               uuid.NewV4().String(),
			FeedingSessionNumber: menuPack.Number,
			MenuItemNumber:       menuItemQuantity.MenuItemNumber,
			Quantity:             menuItemQuantity.Quantity,
			UserNumber:           currentUser,
		}

		err = cateringRepositories.SaveFeedingSessionItem(saveItemAssembly, tx)
		if err != nil {
			return errors.New("saving menu pack item error: " + err.Error())
		}

	}

	return nil
}

func checkDuplicates(menuItemAndQuantities []requests.MenuItemQuantity) error {

	cleanClassNumbers := make([]requests.MenuItemQuantity, len(menuItemAndQuantities))

	for _, selectedMenuItem := range menuItemAndQuantities {

		for _, cleanClassNumber := range cleanClassNumbers {

			if cleanClassNumber.MenuItemNumber == selectedMenuItem.MenuItemNumber {
				menuItem, err := cateringRepositories.GetMenuItemByNumber(selectedMenuItem.MenuItemNumber)
				if err != nil {
					return errors.New("error fetching item details: " + err.Error())
				}
				return errors.New(fmt.Sprintf("%s is selected multiple times", menuItem.Name))
			}

		}

		cleanClassNumbers = append(cleanClassNumbers, selectedMenuItem)

	}

	return nil
}

func GetFeedingSession(schoolNumber string) ([]responses.FeedingSessionResponse, error) {

	var feedingSessions []responses.FeedingSessionResponse

	sessions, err := cateringRepositories.GetFeedingSessionsBySchoolNumber(schoolNumber)

	if err != nil {
		return feedingSessions, err
	}

	for _, session := range sessions {
		feedingSessions = append(feedingSessions, BuildFeedingSessionResponse(session))
	}

	return feedingSessions, nil
}

func BuildFeedingSessionResponse(session models.FeedingSession) responses.FeedingSessionResponse {

	response := responses.FeedingSessionResponse{
		Name:                 session.Name,
		SchoolNumber:         session.SchoolNumber,
		FeedingSessionNumber: session.Number,
	}

	user, err := repository.GetUserByNumber(session.UserNumber)
	if err != nil {
		fmt.Println("Getting session user error", err.Error())
	} else if user.Number == "" {
		response.AddedBy = user.FirstName
	}

	var packItems = make([]responses.PackItem, 0)

	var packCost float64

	sessionItems, err := cateringRepositories.GetFeedingSessionItems(session.Number)
	for _, sessionItem := range sessionItems {
		menuItem, _ := cateringRepositories.GetMenuItemByNumber(sessionItem.MenuItemNumber)
		unitOfUse, _ := cateringRepositories.GetUnitOfMeasureByNumber(menuItem.UnitOfMeasureNumber)
		unitCost := determineUnitCost(menuItem)
		packCost = packCost + commons.RoundFloat(unitCost, 9)*sessionItem.Quantity
		packItems = append(packItems, responses.PackItem{
			ItemNumber: menuItem.Number,
			Name:       menuItem.Name,
			Quantity:   sessionItem.Quantity,
			Type:       menuItem.Type,
			Price:      commons.RoundFloat(unitCost, 6),
			UnitOfConsumption: responses.UnitOfMeasureResponse{
				UnitOfMeasureNumber: unitOfUse.Number,
				Name:                unitOfUse.Name,
				Abbreviation:        unitOfUse.Abbreviation,
				Code:                unitOfUse.Code,
			},
		})
	}

	response.PackItems = packItems
	response.PackCost = commons.RoundFloat(packCost, 6)

	return response
}

func MealUsers(schoolNumber string, campusNumber string, groupNumber string) ([]responses.FeedingUserGroupResponse, error) {

	var response = make([]responses.FeedingUserGroupResponse, 0)

	users, err := cateringRepositories.GetFeedingGroupUserByGroupNumber(schoolNumber, campusNumber, groupNumber)
	if err != nil {
		return nil, errors.New("Getting group users error: " + err.Error())
	}

	var userNumbers []string
	for _, u := range users {
		if u.User.Type == "Student" {
			userNumbers = append(userNumbers, u.UserNumber)
		}
	}

	enrollmentInfoMap := make(map[string]responses.FeedingUserGroupResponse)

	if len(userNumbers) > 0 {
		var enrollmentInfos []struct {
			UserNumber   string
			Adn          string
			Class        string
			Level        string
			AcademicYear string
			StudyPeriod  string
		}

		err := config.DB.Table("students").
			Select("students.user_number as user_number, students.adn as adn, classes.name as class, levels.name as level, academic_years.academic_year_name as academic_year, study_period_names.name as study_period").
			Joins("LEFT JOIN enrollments ON enrollments.student_number = students.number AND enrollments.status = ?", constants.StatusActive).
			Joins("LEFT JOIN campus_classes ON campus_classes.number = enrollments.campus_class_number").
			Joins("LEFT JOIN curriculum_classes ON curriculum_classes.number = campus_classes.curriculum_class_number").
			Joins("LEFT JOIN classes ON classes.number = curriculum_classes.class_number").
			Joins("LEFT JOIN levels ON levels.number = classes.level_number").
			Joins("LEFT JOIN study_periods ON study_periods.number = enrollments.study_period_number").
			Joins("LEFT JOIN study_period_names ON study_period_names.number = study_periods.study_period_name_number").
			Joins("LEFT JOIN academic_years ON academic_years.number = study_periods.academic_year_number").
			Where("students.user_number IN ?", userNumbers).
			Scan(&enrollmentInfos).Error

		if err != nil {
			fmt.Println("Error fetching enrollment info in MealUsers:", err)
		} else {
			for _, info := range enrollmentInfos {
				enrollmentInfoMap[info.UserNumber] = responses.FeedingUserGroupResponse{
					ADN:          info.Adn,
					Class:        info.Class,
					Level:        info.Level,
					AcademicYear: info.AcademicYear,
					StudyPeriod:  info.StudyPeriod,
				}
			}
		}
	}

	for _, user := range users {
		res := createFeedingUserResponse(user)
		if info, ok := enrollmentInfoMap[user.UserNumber]; ok {
			res.ADN = info.ADN
			res.Class = info.Class
			res.Level = info.Level
			res.AcademicYear = info.AcademicYear
			res.StudyPeriod = info.StudyPeriod
		}
		response = append(response, res)
	}

	return response, nil
}

func createFeedingUserResponse(user models.FeedingGroupUser) responses.FeedingUserGroupResponse {

	_user := user.User

	userNames := _user.FirstName + " " + _user.LastName

	return responses.FeedingUserGroupResponse{
		Number:             user.Number,
		UserNumber:         user.UserNumber,
		FeedingGroupName:   user.FeedingGroup.Name,
		FeedingGroupNumber: user.FeedingGroupNumber,
		CampusName:         user.Campus.Name,
		Gender:             user.User.Gender,
		Name:               userNames,
	}
}

func AvailableUser(schoolNumber string, userType string) ([]authResponse.UserResponse, error) {

	var response = make([]authResponse.UserResponse, 0)

	users, err := cateringRepositories.AvailableUsers(schoolNumber, constants.StatusActive, userType)
	if err != nil {
		return nil, errors.New("Getting available users error: " + err.Error())
	}

	for _, user := range users {
		response = append(response, services.CreateUserResponse(user))
	}

	return response, nil
}

func FeedingGroupLearners(schoolNumber, campusNumber, groupNumber, academicYearNumber, studyPeriodNumber, classNumber, levelNumber string) ([]responses.FeedingUserGroupResponse, error) {
	db := config.DB

	var response = make([]responses.FeedingUserGroupResponse, 0)

	var results []struct {
		UserNumber         string
		Number             string
		Gender             string
		FirstName          string
		LastName           string
		Surname            string
		SecondName         string
		OtherName          string
		Photo              string
		FeedingGroupNumber string
		FeedingGroupName   string
		CampusNumber       string
		CampusName         string
		SchoolNumber       string
		SchoolName         string
		Adn                string
		Class              string
		Level              string
		AcademicYear       string
		StudyPeriod        string
		EnrollmentNumber   string
		EnrollmentStatus   string
		EnrollmentDate     string
		StudyPeriodStart   string
		StudyPeriodEnd     string
		ContinuingStudent  bool
		CampusClassNumber  string
		StudyPeriodNumber  string
		StreamNumber       string
	}

	query := db.Table("feeding_group_users").
		Select("feeding_group_users.user_number, feeding_group_users.number, users.gender, users.first_name, users.last_name, " +
			"students.surname, students.second_name, students.other_name, students.photo, " +
			"feeding_group_users.feeding_group_number, feeding_groups.name as feeding_group_name, " +
			"feeding_group_users.campus_number as campus_number, campuses.name as campus_name, " +
			"feeding_groups.school_number as school_number, schools.name as school_name, " +
			"students.adn as adn, classes.name as class, levels.name as level, " +
			"academic_years.academic_year_name as academic_year, study_period_names.name as study_period, " +
			"enrollments.number as enrollment_number, enrollments.status as enrollment_status, " +
			"enrollments.enrollment_date as enrollment_date, " +
			"study_periods.study_period_start_date as study_period_start, " +
			"study_periods.study_period_end_date as study_period_end, " +
			"enrollments.continuing_student as continuing_student, " +
			"enrollments.campus_class_number as campus_class_number, " +
			"enrollments.study_period_number as study_period_number, " +
			"enrollments.stream_number as stream_number").
		Joins("JOIN users ON feeding_group_users.user_number = users.number").
		Joins("JOIN feeding_groups ON feeding_group_users.feeding_group_number = feeding_groups.number").
		Joins("JOIN campuses ON feeding_group_users.campus_number = campuses.number").
		Joins("JOIN schools ON feeding_groups.school_number = schools.number").
		Joins("JOIN students ON students.user_number = users.number").
		Joins("LEFT JOIN enrollments ON enrollments.student_number = students.number AND enrollments.status = 'ACTIVE'").
		Joins("LEFT JOIN campus_classes ON campus_classes.number = enrollments.campus_class_number").
		Joins("LEFT JOIN curriculum_classes ON curriculum_classes.number = campus_classes.curriculum_class_number").
		Joins("LEFT JOIN classes ON classes.number = curriculum_classes.class_number").
		Joins("LEFT JOIN levels ON levels.number = classes.level_number").
		Joins("LEFT JOIN study_periods ON study_periods.number = enrollments.study_period_number").
		Joins("LEFT JOIN study_period_names ON study_period_names.number = study_periods.study_period_name_number").
		Joins("LEFT JOIN academic_years ON academic_years.number = study_periods.academic_year_number").
		Where("users.type = ?", "Student")

	if schoolNumber != "" {
		query = query.Where("feeding_groups.school_number = ?", schoolNumber)
	}
	if campusNumber != "" {
		query = query.Where("feeding_group_users.campus_number = ?", campusNumber)
	}
	if groupNumber != "" {
		query = query.Where("feeding_group_users.feeding_group_number = ?", groupNumber)
	}
	if academicYearNumber != "" {
		query = query.Where("study_periods.academic_year_number = ?", academicYearNumber)
	}
	if studyPeriodNumber != "" {
		query = query.Where("enrollments.study_period_number = ?", studyPeriodNumber)
	}
	if classNumber != "" {
		query = query.Where("curriculum_classes.class_number = ?", classNumber)
	}
	if levelNumber != "" {
		query = query.Where("classes.level_number = ?", levelNumber)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	for _, r := range results {
		response = append(response, responses.FeedingUserGroupResponse{
			// Feeding group
			FeedingGroupNumber: r.FeedingGroupNumber,
			FeedingGroupName:   r.FeedingGroupName,
			// User / student identity
			UserNumber:        r.UserNumber,
			Number:            r.Number,
			Gender:            r.Gender,
			Name:              r.FirstName + " " + r.LastName,
			StudentSurname:    r.Surname,
			StudentSecondName: r.SecondName,
			StudentOtherName:  r.OtherName,
			ADN:               r.Adn,
			Photo:             commons.GenerateDownloadPath(r.Photo, true),
			// Enrollment
			EnrollmentNumber:     r.EnrollmentNumber,
			EnrollmentStatus:     r.EnrollmentStatus,
			EnrollmentDate:       r.EnrollmentDate,
			StudyPeriodStartDate: r.StudyPeriodStart,
			StudyPeriodEndDate:   r.StudyPeriodEnd,
			ContinuingStudent:    r.ContinuingStudent,
			// Location / class
			SchoolNumber:      r.SchoolNumber,
			School:            r.SchoolName,
			CampusNumber:      r.CampusNumber,
			CampusName:        r.CampusName,
			CampusClassNumber: r.CampusClassNumber,
			StudyPeriodNumber: r.StudyPeriodNumber,
			StreamNumber:      r.StreamNumber,
			Class:             r.Class,
			Level:             r.Level,
			AcademicYear:      r.AcademicYear,
			StudyPeriod:       r.StudyPeriod,
		})
	}

	return response, nil
}

func DeleteFeedingGroup(feedingGroupNumber string) error {

	group, err := cateringRepositories.GetFeedingGroupByNumber(feedingGroupNumber)
	if err != nil {
		return errors.New("Checking group number error: " + err.Error())
	}

	if group.Number == "" {
		return errors.New("invalid group number")
	}

	err = cateringRepositories.RemoveFeedingGroup(group)
	if err != nil {
		fmt.Println("Removing feed group error: " + err.Error())
		return errors.New("Can not remove '" + group.Name + "' group, try removing all its usages first.")
	}

	return nil

}

func UpdateFeedingSession(updateFeedingSessionRequest requests.UpdateFeedingSessionRequest, currentUserNumber string) error {
	feedingSession, err := cateringRepositories.GetFeedingSessionByNumberAndSchoolNumber(updateFeedingSessionRequest.FeedingSessionNumber, updateFeedingSessionRequest.SchoolNumber)

	if err != nil {
		return errors.New("Checking feed session number error: " + err.Error())
	}

	if feedingSession.Number == "" {
		return errors.New("invalid feed session number")
	}

	checkName, err := cateringRepositories.GetFeedingSessionByNameAndSchool(updateFeedingSessionRequest.Name, feedingSession.SchoolNumber)
	if err != nil {
		return errors.New("Checking feed session name error: " + err.Error())
	}

	if checkName.Number != "" && checkName.Number != feedingSession.Number {
		return errors.New("name already used, choose a unique one")
	}

	menuPackItems, err := cateringRepositories.GetFeedingSessionItems(feedingSession.Number)
	if err != nil {
		return errors.New("checking old items error: " + err.Error())
	}

	feedingSession.UserNumber = currentUserNumber
	feedingSession.Name = updateFeedingSessionRequest.Name
	feedingSession.Status = updateFeedingSessionRequest.Status

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		for _, menuPackItem := range menuPackItems {
			err = cateringRepositories.RemoveFeedingSessionItem(menuPackItem, tx)
			if err != nil {
				return errors.New("clearing old item error: " + err.Error())
			}
		}

		err = cateringRepositories.SaveFeedingSession(feedingSession, tx)
		if err != nil {
			return errors.New("updating feed session error: " + err.Error())
		}

		err = processMenuItems(feedingSession, updateFeedingSessionRequest.MenuItems, currentUserNumber, tx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func DeleteFeedingSession(feedingSessionNumber string) error {

	feedingSession, err := cateringRepositories.GetFeedingSessionByNumberAndSchoolNumber(feedingSessionNumber, "")
	if err != nil {
		return errors.New("Checking feed session number error: " + err.Error())
	}

	if feedingSession.Number == "" {
		return errors.New("invalid feed session number")
	}

	err = cateringRepositories.DeleteFeedingSession(feedingSession)
	if err != nil {
		fmt.Println("Deleting feed session error: " + err.Error())
		return errors.New("this feeding session could not be removed, try removing all its usages first")
	}

	return nil
}

func SaveFeedingGroupCalendar(feedingGroupCalendarRequest requests.FeedingGroupCalendarRequest, currentUser string) error {

	trn := config.DB.Begin()

	for _, groupNumber := range feedingGroupCalendarRequest.FeedingGroups {
		err := processGroupCalendar(groupNumber, trn, feedingGroupCalendarRequest.CalendarEvents, currentUser)
		if err != nil {
			trn.Rollback()
			return err
		}
	}

	if err := trn.Commit().Error; err != nil {
		return errors.New(fmt.Sprintf("Finalizing groups calendars error: %v", err))
	}

	return nil
}

func processGroupCalendar(groupNumber string, trn *gorm.DB, calendarEvents []string, currentUser string) error {

	groupCalendars, err := cateringRepositories.FindCalendarEventFeedingGroupsByGroupNumber(groupNumber, trn)
	if err != nil {
		return errors.New("Finding group calendars error: " + err.Error())
	}

	for _, groupCalendar := range groupCalendars {
		err = cateringRepositories.RemoveCalendarEventFeedingGroup(groupCalendar, trn)
		if err != nil {
			return errors.New("Clearing previous group calendar error: " + err.Error())
		}
	}

	for _, calendarEvent := range calendarEvents {

		checkDuplicate, err := cateringRepositories.GetCalendarEventFeedingGroupsByGroupNumber(groupNumber, calendarEvent, trn)
		if err != nil {
			return errors.New("Checking group calendar duplicates error: " + err.Error())
		}

		if checkDuplicate.Number != "" {
			continue
		}

		err = cateringRepositories.SavedCalendarEventFeedingGroup(models.CalendarEventFeedingGroup{
			Number:              uuid.NewV4().String(),
			UserNumber:          currentUser,
			FeedingGroupNumber:  groupNumber,
			CalendarEventNumber: calendarEvent,
		}, trn)

		if err != nil {
			return errors.New("Activating group calendars error: " + err.Error())
		}
	}

	return nil
}

func GetFeedingGroupCalendarEvent(schoolNumber string, entityNumber string) (responses.FeedingGroupCalendarEvent, error) {

	var feedingGroupCalendarEvent responses.FeedingGroupCalendarEvent

	endDateAfter := time.Now()

	calendarEvents, err := repositories.FindCalendarEvent(entityNumber, endDateAfter)
	if err != nil {
		return feedingGroupCalendarEvent, errors.New("Getting calendar event error: " + err.Error())
	}

	var feedingGroupCalendarEventItems []responses.FeedingGroupCalendarEventItem
	feedingGroupCalendarEventItems = append(feedingGroupCalendarEventItems, responses.FeedingGroupCalendarEventItem{})

	for _, calendarEvent := range calendarEvents {
		feedingGroupCalendarEventItems = append(feedingGroupCalendarEventItems, calendarEventFrom(calendarEvent))
	}

	feedingGroupCalendarEvent.CalendarEvents = feedingGroupCalendarEventItems

	feedingGroups, err := cateringRepositories.GetAllFeedingGroups(schoolNumber, constants.StatusActive)
	if err != nil {
		return feedingGroupCalendarEvent, errors.New("Getting groups error: " + err.Error())
	}

	var calendarGroupItems []responses.CalendarGroupItem

	for _, feedingGroup := range feedingGroups {
		calendarGroupItem, err := createGroupCalendarItem(feedingGroup, calendarEvents)

		if err != nil {
			return feedingGroupCalendarEvent, err
		}

		calendarGroupItems = append(calendarGroupItems, calendarGroupItem)
	}

	feedingGroupCalendarEvent.CalendarGroups = calendarGroupItems

	return feedingGroupCalendarEvent, nil
}

func createGroupCalendarItem(feedingGroup models.FeedingGroup, calendarEvents []eosModels.CalendarEvent) (responses.CalendarGroupItem, error) {

	calendarGroupItem := responses.CalendarGroupItem{
		GroupName: feedingGroup.Name,
		Number:    feedingGroup.Number,
	}

	var activeCalendarItems []responses.ActiveCalendarItem

	for _, calendarEvent := range calendarEvents {
		activeCalendarItem := responses.ActiveCalendarItem{}
		activeCalendarItem.FeedingGroupCalendarEventItem = calendarEventFrom(calendarEvent)

		check, err := cateringRepositories.GetCalendarEventFeedingGroupsByGroupNumber(feedingGroup.Number, calendarEvent.Number)
		if err != nil {
			return calendarGroupItem, errors.New("Checking group " + feedingGroup.Number + " on calendar " + calendarEvent.Name + " error: " + err.Error())
		}

		if check.Number == "" {
			activeCalendarItem.Active = false
			activeCalendarItem.Message = "Not active during this period"
		} else {
			activeCalendarItem.Active = true
			activeCalendarItem.Message = "Activated for this period"
		}

		activeCalendarItems = append(activeCalendarItems, activeCalendarItem)
	}

	calendarGroupItem.ActiveCalendars = activeCalendarItems

	return calendarGroupItem, nil
}

func calendarEventFrom(calendarEvent eosModels.CalendarEvent) responses.FeedingGroupCalendarEventItem {
	return responses.FeedingGroupCalendarEventItem{
		CalendarName: calendarEvent.Name,
		Number:       calendarEvent.Number,
		StartDate:    calendarEvent.StartDate.Format(time.DateTime),
		EndDate:      calendarEvent.EndDate.Format(time.DateTime),
	}
}

func UserMeals(userNumber string) (responses.CalendarGroupItem, error) {

	userGroup, err := cateringRepositories.GetFeedingGroupUserByUserNumber(userNumber)
	if err != nil {
		return responses.CalendarGroupItem{}, errors.New("Fetching user group error: " + err.Error())
	}

	if userGroup.Number == "" {
		return responses.CalendarGroupItem{}, errors.New("not attached to any group for meals")
	}

	group, err := cateringRepositories.GetFeedingGroupByNumber(userGroup.FeedingGroupNumber)
	if err != nil {
		return responses.CalendarGroupItem{}, errors.New("Getting group details error: " + err.Error())
	}

	school, err := schoolRepositories.GetSchoolByNumber(group.SchoolNumber)
	if err != nil {
		return responses.CalendarGroupItem{}, errors.New("Checking school error: " + err.Error())
	}

	endDateAfter := time.Now()

	calendarEvents, err := repositories.FindCalendarEvent(school.EntityNumber, endDateAfter)
	if err != nil {
		return responses.CalendarGroupItem{}, errors.New("Getting group details error: " + err.Error())
	}

	userMealResponse, err := createGroupCalendarItem(group, calendarEvents)
	if err != nil {
		return responses.CalendarGroupItem{}, err
	}

	return userMealResponse, nil
}
