package controllers

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	campusModes "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	subjectModes "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func AddPeriod(context *gin.Context) {
	var timetable models.StudyTimetablePeriods
	var timetableRequest requests.TimetableRequest

	err := context.ShouldBind(&timetableRequest)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var timetableRequestDayCode int32
	// assign day code	1 to monday , 2 to tuesday and so on
	switch timetableRequest.Day {
	case "Monday":
		timetableRequestDayCode = 1
	case "Tuesday":
		timetableRequestDayCode = 2
	case "Wednesday":
		timetableRequestDayCode = 3
	case "Thursday":
		timetableRequestDayCode = 4
	case "Friday":
		timetableRequestDayCode = 5
	case "Saturday":
		timetableRequestDayCode = 6
	case "Sunday":
		timetableRequestDayCode = 7
	default:
		timetableRequestDayCode = 0

	}

	timetable.Number = uuid.NewV4().String()

	period := models.StudyTimetablePeriods{
		Number:        timetable.Number,
		StartTime:     timetableRequest.StartTime,
		EndTime:       timetableRequest.EndTime,
		SubjectNumber: timetableRequest.SubjectNumber,
		StreamNumber:  timetableRequest.StreamNumber,
		Day:           timetableRequest.Day,
		DayCode:       timetableRequestDayCode,
		CampusNumber:  timetableRequest.CampusNumber,
		UserNumber:    context.GetHeader("User-Number"),
	}

	config.DB.Save(&period)

	context.JSON(http.StatusOK, gin.H{"message": "Period added successfully"})

}

func GetDistinctDays(ctx *gin.Context) {
	var days []string
	config.DB.Model(&models.StudyTimetablePeriods{}).
		Select("DISTINCT day").
		Order("day_code asc").
		Scan(&days)
	ctx.JSON(http.StatusOK, gin.H{"days": days})

}

func GetPeriodsByDay(ctx *gin.Context) {
	day := ctx.Param("day")
	var periods []models.StudyTimetablePeriods
	config.DB.Where("day = ?", day).Find(&periods)
	ctx.JSON(http.StatusOK, gin.H{"data": periods})
}

func FetchTimetablePeriodsGroupedByDay(ctx *gin.Context) {
	var periods []models.StudyTimetablePeriods
	config.DB.
		//Model(&models.StudyTimetablePeriods{}).
		//Select("stream_number,day").
		Preload("Subject").
		Preload("Stream").
		Group("day_code").
		Order("day_code asc").
		Joins("left join subjects on subjects.number=timetables.subject_number").
		Find(&periods)

	ctx.JSON(http.StatusOK, gin.H{"data": periods})
}

func GetPeriodsByStream(ctx *gin.Context) {
	streamNumber := ctx.Param("stream_number")
	var periods []models.StudyTimetablePeriods
	config.DB.Where("stream_number = ?", streamNumber).Find(&periods)
	ctx.JSON(http.StatusOK, gin.H{"data": periods})
}

func FetchClassLevelsByCampusNumberAndLevelNumber(ctx *gin.Context) {

	campusNumber := ctx.Query("campus_number")
	levelNumber := ctx.Query("level_number")

	//fmt.Println(campusNumber, levelNumber)

	var classLevels []classModels.ClassLevel

	config.DB.
		Where("campus_number = ? AND level_number = ?", campusNumber, levelNumber).
		Preload("Class").
		Find(&classLevels)
	ctx.JSON(http.StatusOK, gin.H{"data": classLevels})

}

func FetchSubjectsByClassLevelNumber(ctx *gin.Context) {
	classLevelNumber := ctx.Param("class_level_number")
	var subjects []responses.Subjects

	config.DB.Model(subjectModes.Subject{}).
		Select("subjects.name,subjects.number").
		Joins("left join class_level_subjects on class_level_subjects.subject_number=subjects.number").
		Where("class_level_subjects.class_level_number = ?", classLevelNumber).
		Scan(&subjects)
	ctx.JSON(http.StatusOK, gin.H{"data": subjects})
}

func FetchLevelsByCampusNumber(ctx *gin.Context) {
	campusNumber := ctx.Param("campus_number")
	var levels []campusModes.Level

	config.DB.
		Where("campus_number = ?", campusNumber).
		Find(&levels)

	ctx.JSON(http.StatusOK, gin.H{"data": levels})
}

/*func AddCalendarEvent(ctx *gin.Context) {

	var calendarEvent eosModels.CalendarEvent
	var calendarEventRequest requests.CalendarEventRequest

	err := ctx.BindJSON(&calendarEventRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	calendarEvent.Number = uuid.NewV4().String()

	//convert date to time
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, calendarEventRequest.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	endDate, _ := time.Parse(layout, calendarEventRequest.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := eosModels.CalendarEvent{
		Number:    calendarEvent.Number,
		StartDate: startDate,
		EndDate:   endDate,
		StartTime: calendarEventRequest.StartTime,
		EndTime:   calendarEventRequest.EndTime,
	}

	config.DB.Save(&event)

	ctx.JSON(http.StatusOK, gin.H{"message": "Event added successfully"})
}

func CreateGeneralCalendar(ctx *gin.Context) {
	var generalCalendar eosModels.GeneralCalendar
	var generalCalendarRequest requests.GeneralCalendarRequest

	err := ctx.BindJSON(&generalCalendarRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	generalCalendar.Number = uuid.NewV4().String()

	//convert date to time
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, generalCalendarRequest.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	endDate, _ := time.Parse(layout, generalCalendarRequest.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := eosModels.GeneralCalendar{
		Number:    generalCalendar.Number,
		StartDate: startDate,
		EndDate:   endDate,
		Name:      generalCalendarRequest.Name,
	}

	config.DB.Save(&event)

	ctx.JSON(http.StatusOK, gin.H{"message": "Event added successfully"})
}*/

func CreateStudyTimetable(ctx *gin.Context) {

	var request = requests.StudyTimetableRequest{}
	var StudyTimetable = models.StudyTimetable{}

	err := ctx.Bind(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback("Failed to unwrap request"))
		return
	}

	StudyTimetable.Number = uuid.NewV4().String()
	StudyTimetable.UserNumber = ctx.GetHeader("User-Number")
	StudyTimetable.Status = constants.STATUS_PENDING
	StudyTimetable.StudyPeriodNumber = request.StudyPeriodNumber
	StudyTimetable.Name = request.Name
	err = config.DB.Create(&StudyTimetable).Error

	if err != nil {
		fmt.Println("Error", err.Error())
		ctx.JSON(http.StatusBadRequest, commons.Feedback("Failed to save timetable"))
		return
	}

	ctx.JSON(http.StatusOK, commons.Feedback("Timetable added successfully"))
}

func FetchStudyTimetables(ctx *gin.Context) {

	var timetables []models.StudyTimetable
	studyTimetables := []responses.StudyTimetables{}

	err := config.DB.Find(&timetables).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, commons.Feedback("Failed to fetch timetables"))
		return
	}

	for _, timetable := range timetables {
		studyTimetables = append(studyTimetables, responses.StudyTimetables{
			Number: timetable.Number,
			Name:   timetable.Name,
			Status: timetable.Status,
		})
	}

	ctx.JSON(http.StatusOK, studyTimetables)
}
