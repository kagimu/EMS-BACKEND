package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	campusModes "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	subjectModes "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/service"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateStudyTimetablePeriod(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.CreateStudyTimetablePeriodRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	periodService := getStudyTimetablePeriodService()
	if err := periodService.CreateStudyTimetablePeriod(payload, user.Number); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Study timetable period created successfully"})
}

func UpdateStudyTimetablePeriod(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	periodNumber := c.Param("period_number")
	if periodNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Period number is required"})
		return
	}

	var payload requests.UpdateStudyTimetablePeriodRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	periodService := getStudyTimetablePeriodService()
	if err := periodService.UpdateStudyTimetablePeriod(periodNumber, payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Study timetable period updated successfully"})
}

func DeleteStudyTimetablePeriod(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	periodNumber := c.Param("period_number")
	if periodNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Period number is required"})
		return
	}

	periodService := getStudyTimetablePeriodService()
	if err := periodService.DeleteStudyTimetablePeriod(periodNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Study timetable period deleted successfully"})
}

func GetStudyTimetablePeriod(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	periodNumber := c.Param("period_number")
	if periodNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Period number is required"})
		return
	}

	periodService := getStudyTimetablePeriodService()
	period, err := periodService.GetStudyTimetablePeriodByNumber(periodNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": period})
}

func GetStudyTimetablePeriods(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	userNumber := c.Query("user_number")
	if userNumber == "" {
		userNumber = user.Number
	}

	studyTimetableNumber := c.Query("study_timetable_number")
	streamNumber := c.Query("stream_number")
	day := c.Query("day")

	periodService := getStudyTimetablePeriodService()
	periods, err := periodService.ListStudyTimetablePeriods(studyTimetableNumber, streamNumber, day, userNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": periods})
}

func getStudyTimetablePeriodService() service.StudyTimetablePeriodService {
	periodRepo := repository.NewStudyTimetablePeriodRepository(config.DB)
	slotRepo := repository.NewTimetableSlotRepository(config.DB)
	taskRepo := repository.NewTimetableTaskRepository(config.DB)
	timetableRepo := repository.NewStudyTimetableRepository(config.DB)
	return service.NewStudyTimetablePeriodService(periodRepo, slotRepo, taskRepo, timetableRepo)
}

func AddPeriod(context *gin.Context) {
	var timetable models.StudyTimetablePeriod
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

	period := models.StudyTimetablePeriod{
		Number:       timetable.Number,
		StartTime:    timetableRequest.StartTime,
		EndTime:      timetableRequest.EndTime,
		StreamNumber: timetableRequest.StreamNumber,
		Day:          timetableRequest.Day,
		DayCode:      timetableRequestDayCode,
		UserNumber:   context.GetHeader("User-Number"),
	}

	config.DB.Save(&period)

	context.JSON(http.StatusOK, gin.H{"message": "Period added successfully"})

}

func GetDistinctDays(ctx *gin.Context) {
	var days []string
	config.DB.Model(&models.StudyTimetablePeriod{}).
		Select("DISTINCT day").
		Order("day_code asc").
		Scan(&days)
	ctx.JSON(http.StatusOK, gin.H{"days": days})

}

func GetPeriodsByDay(ctx *gin.Context) {
	day := ctx.Param("day")
	var periods []models.StudyTimetablePeriod
	config.DB.Where("day = ?", day).Find(&periods)
	ctx.JSON(http.StatusOK, gin.H{"data": periods})
}

func FetchTimetablePeriodsGroupedByDay(ctx *gin.Context) {
	var periods []models.StudyTimetablePeriod
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
	var periods []models.StudyTimetablePeriod
	config.DB.Where("stream_number = ?", streamNumber).Find(&periods)
	ctx.JSON(http.StatusOK, gin.H{"data": periods})
}

func FetchClassLevelsByCampusNumberAndLevelNumber(ctx *gin.Context) {

	campusNumber := ctx.Query("campus_number")
	levelNumber := ctx.Query("level_number")

	//fmt.Println(campusNumber, levelNumber)

	var classLevels []classModels.CampusClass

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

func CreateTimetableTask(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.TimetableTaskRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	timetableTaskRepo := repository.NewTimetableTaskRepository(config.DB)
	timetableTaskService := service.NewTimetableTaskService(timetableTaskRepo)

	if err := timetableTaskService.CreateTimetableTask(payload, user.Number); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Timetable task created successfully"})
}

func UpdateTimetableTask(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	taskNumber := c.Param("task_number")
	if taskNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Task number is required"})
		return
	}

	var payload requests.UpdateTimetableTaskRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	timetableTaskRepo := repository.NewTimetableTaskRepository(config.DB)
	timetableTaskService := service.NewTimetableTaskService(timetableTaskRepo)

	if err := timetableTaskService.UpdateTimetableTask(taskNumber, payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Timetable task updated successfully"})
}

func DeleteTimetableTask(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	taskNumber := c.Param("task_number")
	if taskNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Task number is required"})
		return
	}

	timetableTaskRepo := repository.NewTimetableTaskRepository(config.DB)
	timetableTaskService := service.NewTimetableTaskService(timetableTaskRepo)

	if err := timetableTaskService.DeleteTimetableTask(taskNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Timetable task deleted successfully"})
}

func GetTimetableTasks(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	classNumber := c.Query("class_number")

	timetableTaskRepo := repository.NewTimetableTaskRepository(config.DB)
	timetableTaskService := service.NewTimetableTaskService(timetableTaskRepo)

	tasks, err := timetableTaskService.GetAllTimetableTasks(classNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func GetTimetableTask(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	taskNumber := c.Param("task_number")
	if taskNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Task number is required"})
		return
	}

	timetableTaskRepo := repository.NewTimetableTaskRepository(config.DB)
	timetableTaskService := service.NewTimetableTaskService(timetableTaskRepo)

	task, err := timetableTaskService.GetTimetableTaskByNumber(taskNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

func CreateTimetableSlot(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.TimetableSlotRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	timetableSlotRepo := repository.NewTimetableSlotRepository(config.DB)
	timetableSlotService := service.NewTimetableSlotService(timetableSlotRepo)

	if err := timetableSlotService.CreateTimetableSlot(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Timetable slot created successfully"})
}

func UpdateTimetableSlot(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	slotNumber := c.Param("slot_number")
	if slotNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Slot number is required"})
		return
	}

	var payload requests.UpdateTimetableSlotRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	timetableSlotRepo := repository.NewTimetableSlotRepository(config.DB)
	timetableSlotService := service.NewTimetableSlotService(timetableSlotRepo)

	if err := timetableSlotService.UpdateTimetableSlot(slotNumber, payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Timetable slot updated successfully"})
}

func DeleteTimetableSlot(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	slotNumber := c.Param("slot_number")
	if slotNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Slot number is required"})
		return
	}

	timetableSlotRepo := repository.NewTimetableSlotRepository(config.DB)
	timetableSlotService := service.NewTimetableSlotService(timetableSlotRepo)

	if err := timetableSlotService.DeleteTimetableSlot(slotNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Timetable slot deleted successfully"})
}

func GetTimetableSlots(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	schoolNumber := c.Query("school_number")
	levelNumber := c.Query("level_number")

	timetableSlotRepo := repository.NewTimetableSlotRepository(config.DB)
	timetableSlotService := service.NewTimetableSlotService(timetableSlotRepo)

	slots, err := timetableSlotService.GetAllTimetableSlots(schoolNumber, levelNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": slots})
}

func GetTimetableSlot(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	slotNumber := c.Param("slot_number")
	if slotNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Slot number is required"})
		return
	}

	timetableSlotRepo := repository.NewTimetableSlotRepository(config.DB)
	timetableSlotService := service.NewTimetableSlotService(timetableSlotRepo)

	slot, err := timetableSlotService.GetTimetableSlotByNumber(slotNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": slot})
}
