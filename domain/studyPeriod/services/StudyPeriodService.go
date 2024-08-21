package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/api/responses"
	studyModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type StudyPeriodService interface {
	StudyPeriodDetails(studyPeriodNumber string) (responses.StudyPeriodListDetail, error)
	NewStudyPeriod(studyPeriodRequest requests.StudyPeriodRequest, currentUserNumber string) error
	NewStudyPeriodName(studyPeriodNameRequest requests.StudyPeriodNameRequest, currentUserNumber string) error
	UpdateStudyPeriod(studyPeriodRequest requests.UpdateStudyPeriodRequest, currentUserNumber string) error
	UpdateStudyPeriodName(studyPeriodNameRequest requests.UpdateStudyPeriodNameRequest, currentUserNumber string) error
	FindAllStudyPeriodNames(schoolNumber string) ([]responses.StudyPeriodNameResponse, error)
	FindAllStudyPeriods(academicYearNumber, calendarEventNumber, studyPeriodNameNumber, classNumber string) ([]responses.StudyPeriodListDetail, error)
	RemoveStudyPeriodName(studyPeriodNameNumber string) error
	RemoveStudyPeriod(studyPeriodNumber string) error
}

type studyPeriodService struct {
}

func (studyPeriodService *studyPeriodService) FindAllStudyPeriods(academicYearNumber, calendarEventNumber, studyPeriodNameNumber, classNumber string) ([]responses.StudyPeriodListDetail, error) {
	var studyPeriodsResponse = make([]responses.StudyPeriodListDetail, 0)

	studyPeriods, err := repositories.FindStudyPeriod(academicYearNumber, calendarEventNumber, studyPeriodNameNumber, classNumber)
	if err != nil {
		return nil, errors.New("Filtering study period error " + err.Error())
	}

	for _, studyPeriod := range studyPeriods {
		details, err := createStudyPeriodResponse(studyPeriod)
		if err != nil {
			return nil, err
		}

		studyPeriodsResponse = append(studyPeriodsResponse, details)
	}

	return studyPeriodsResponse, nil
}

func (studyPeriodService *studyPeriodService) StudyPeriodDetails(studyPeriodNumber string) (responses.StudyPeriodListDetail, error) {

	var studyPeriodResponse responses.StudyPeriodListDetail

	studyPeriod, err := repositories.GetStudyPeriodByNumber(studyPeriodNumber)
	if err != nil {
		return studyPeriodResponse, errors.New("Checking study period number error: " + err.Error())
	}

	if studyPeriod.Number == "" {
		return studyPeriodResponse, errors.New("Invalid study period number ")
	}

	return createStudyPeriodResponse(studyPeriod)
}

func createStudyPeriodResponse(studyPeriod studyModels.StudyPeriod) (responses.StudyPeriodListDetail, error) {

	classes, err := repositories.FindClassStudyPeriods("", studyPeriod.Number)
	if err != nil {
		return responses.StudyPeriodListDetail{}, errors.New("Getting study period classes error: " + err.Error())
	}

	var classesResponse = make([]responses2.ClassResponse, 0)

	for _, classNumber := range classes {
		class, _ := repositories3.GetClassByNumber(classNumber.ClassNumber)
		classesResponse = append(classesResponse, services.CreateClassResponse(class))
	}

	response := responses.StudyPeriodListDetail{
		Classes: classesResponse,
	}

	response.StudyPeriodList = createStudyPeriod(studyPeriod)

	if len(classesResponse) > 0 {
		response.LevelName = classesResponse[0].LevelName
	}

	return response, nil
}

func createStudyPeriod(studyPeriod studyModels.StudyPeriod) responses.StudyPeriodList {

	academicYear, _ := repositories4.GetAcademicYearByNumber(studyPeriod.AcademicYearNumber)
	studyPeriodName, _ := repositories.GetStudyPeriodNameByNumber(studyPeriod.StudyPeriodNameNumber)

	return responses.StudyPeriodList{
		CampusLevelNumber:     academicYear.CampusLevelNumber,
		StudyPeriodName:       studyPeriodName.Name,
		StudyPeriodNameNumber: studyPeriodName.Number,
		AcademicYearName:      academicYear.AcademicYearName,
		AcademicYearNumber:    studyPeriod.AcademicYearNumber,
		StudyPeriodNumber:     studyPeriod.Number,
		Code:                  studyPeriod.Code,
		StudyPeriodStartDate:  studyPeriod.StudyPeriodStartDate.Format(time.DateTime),
		StudyPeriodEndDate:    studyPeriod.StudyPeriodEndDate.Format(time.DateTime),
		UserNumber:            studyPeriod.UserNumber,
	}
}

func (studyPeriodService *studyPeriodService) NewStudyPeriod(studyPeriodRequest requests.StudyPeriodRequest, currentUserNumber string) error {

	academicYear, err := repositories4.GetAcademicYearByNumber(studyPeriodRequest.AcademicYearNumber)
	if err != nil {
		return errors.New("Checking academic year number error " + err.Error())
	}

	if academicYear.Number == "" {
		return errors.New("Invalid academic year number provided ")
	}

	layout := "2006-01-02T15:04:05.000Z"
	startDate, errorDate := time.Parse(layout, studyPeriodRequest.StudyPeriodStartDate)
	if errorDate != nil {
		return errorDate
	}

	endDate, errorDate := time.Parse(layout, studyPeriodRequest.StudyPeriodEndDate)
	if errorDate != nil {
		return errorDate
	}

	checkDuplicate, err := repositories.GetStudyPeriodByCode(studyPeriodRequest.AcademicYearNumber, studyPeriodRequest.Code)
	if err != nil {
		return errors.New("Checking code error " + err.Error())
	}

	studyPeriodName, _ := repositories.GetStudyPeriodNameByNumber(studyPeriodRequest.StudyPeriodNameNumber)

	if checkDuplicate.Number != "" {
		return errors.New("Study period code already used ")
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		var calendarEvent = eosModels.CalendarEvent{}

		// Format the time as "15:04" (24-hour format)
		startTime := startDate.Format("15:04")
		endTime := endDate.Format("15:04")

		calendarEvent.Number = uuid.NewV4().String()
		calendarEvent.StartDate = startDate
		calendarEvent.EndDate = endDate
		calendarEvent.StartTime = startTime
		calendarEvent.EndTime = endTime
		calendarEvent.CalendarNumber = academicYear.CalendarNumber
		calendarEvent.Name = studyPeriodName.Name

		err := repositories2.SaveCalendarEvent(calendarEvent, tx)
		if err != nil {
			return err
		}

		studyPeriod, err := repositories.SaveStudyPeriod(studyModels.StudyPeriod{
			Code:                  studyPeriodRequest.Code,
			Number:                uuid.NewV4().String(),
			UserNumber:            currentUserNumber,
			AcademicYearNumber:    studyPeriodRequest.AcademicYearNumber,
			CalendarEventNumber:   calendarEvent.Number,
			StudyPeriodStartDate:  startDate,
			StudyPeriodEndDate:    endDate,
			StudyPeriodNameNumber: studyPeriodRequest.StudyPeriodNameNumber,
		}, tx)

		if err != nil {
			return err
		}

		err = prepareClasses(studyPeriodRequest.Classes, studyPeriod, tx)
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

func prepareClasses(classes []string, period studyModels.StudyPeriod, tx *gorm.DB) error {

	//Clear all other classes on this study period
	classStudyPeriods, err := repositories.FindClassStudyPeriods("", period.Number, tx)
	if err != nil {
		return errors.New("Checking classes on study period error: " + err.Error())
	}

	for _, classStudyPeriod := range classStudyPeriods {
		err = repositories.DeleteClassStudyPeriod(classStudyPeriod, tx)
		if err != nil {
			return errors.New("Clearing previous classes on study period error: " + err.Error())
		}
	}

	for _, classNumber := range classes {

		//Clear this class form any other study period
		classStudyPeriods, err = repositories.FindClassStudyPeriods(classNumber, "", tx)
		if err != nil {
			return errors.New("Checking classes on study period error: " + err.Error())
		}

		for _, classStudyPeriod := range classStudyPeriods {
			err = repositories.DeleteClassStudyPeriod(classStudyPeriod, tx)
			if err != nil {
				return errors.New("Clearing previous study period on class error: " + err.Error())
			}
		}

		//Attach class to study period
		err = repositories.SaveClassStudyPeriod(studyModels.ClassStudyPeriod{
			ClassNumber:       classNumber,
			StudyPeriodNumber: period.Number,
		}, tx)

		if err != nil {
			return errors.New("Attaching class to study period error: " + err.Error())
		}

	}

	return nil
}

func (studyPeriodService *studyPeriodService) UpdateStudyPeriod(studyPeriodRequest requests.UpdateStudyPeriodRequest, currentUserNumber string) error {

	academicYear, err := repositories4.GetAcademicYearByNumber(studyPeriodRequest.AcademicYearNumber)
	if err != nil {
		return errors.New("Checking academic year number error " + err.Error())
	}

	if academicYear.Number == "" {
		return errors.New("Invalid academic year number provided ")
	}

	studyPeriod, err := repositories.GetStudyPeriodByNumber(studyPeriodRequest.StudyPeriodNumber)
	if err != nil {
		return errors.New("Checking study period number error: " + err.Error())
	}

	if studyPeriod.Number == "" {
		return errors.New("Invalid study period number ")
	}

	checkDuplicate, err := repositories.GetStudyPeriodByCode(studyPeriodRequest.AcademicYearNumber, studyPeriodRequest.Code)
	if err != nil {
		return errors.New("Checking code error " + err.Error())
	}

	if checkDuplicate.Number != "" && studyPeriod.Number != checkDuplicate.Number {
		return errors.New("Study period code already used ")
	}

	studyPeriodName, _ := repositories.GetStudyPeriodNameByNumber(studyPeriodRequest.StudyPeriodNameNumber)

	layout := "2006-01-02T15:04:05.000Z"
	startDate, errorDate := time.Parse(layout, studyPeriodRequest.StudyPeriodStartDate)
	if errorDate != nil {
		return errorDate
	}

	endDate, errorDate := time.Parse(layout, studyPeriodRequest.StudyPeriodEndDate)
	if errorDate != nil {
		return errorDate
	}

	studyPeriod.UserNumber = currentUserNumber
	studyPeriod.Code = studyPeriodRequest.Code
	studyPeriod.AcademicYearNumber = studyPeriodRequest.AcademicYearNumber
	studyPeriod.StudyPeriodStartDate = startDate
	studyPeriod.StudyPeriodEndDate = endDate
	studyPeriod.StudyPeriodNameNumber = studyPeriodRequest.StudyPeriodNameNumber

	calendarEvent, err := repositories2.GetCalendarEventByNumber(studyPeriod.CalendarEventNumber)
	if err != nil {
		return errors.New("Checking calendar event error: " + err.Error())
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		// Format the time as "15:04" (24-hour format)
		startTime := startDate.Format("15:04")
		endTime := endDate.Format("15:04")
		calendarEvent.StartDate = startDate
		calendarEvent.EndDate = endDate
		calendarEvent.StartTime = startTime
		calendarEvent.EndTime = endTime
		calendarEvent.CalendarNumber = academicYear.CalendarNumber
		calendarEvent.Name = studyPeriodName.Name

		err := repositories2.SaveCalendarEvent(calendarEvent, tx)
		if err != nil {
			return err
		}

		studyPeriod, err := repositories.SaveStudyPeriod(studyPeriod, tx)

		if err != nil {
			return err
		}

		err = prepareClasses(studyPeriodRequest.Classes, studyPeriod, tx)
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

func (studyPeriodService *studyPeriodService) RemoveStudyPeriod(studyPeriodNumber string) error {
	studyPeriod, err := repositories.GetStudyPeriodByNumber(studyPeriodNumber)
	if err != nil {
		return errors.New("Checking study period error: " + err.Error())
	}

	if studyPeriod.Number == "" {
		return errors.New("Invalid study period number ")
	}

	calendarEvent, err := repositories2.GetCalendarEventByNumber(studyPeriod.CalendarEventNumber)
	if err != nil {
		return err
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {

		err = repositories.RemoveStudyPeriod(studyPeriod, tx)
		if err != nil {
			return errors.New("Removing study period error: " + err.Error())
		}

		err = repositories2.RemoveCalendarEvent(calendarEvent, tx)
		if err != nil {
			return errors.New("Removing calendar event error: " + err.Error())
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func NewStudyPeriodService() StudyPeriodService {
	return &studyPeriodService{}
}

func (studyPeriodService *studyPeriodService) NewStudyPeriodName(studyPeriodNameRequest requests.StudyPeriodNameRequest, currentUserNumber string) error {

	checkName, err := repositories.GetStudyPeriodNameByName(studyPeriodNameRequest.SchoolNumber, studyPeriodNameRequest.Name)
	if err != nil {
		return errors.New("Checking name duplicate error " + err.Error())
	}

	if checkName.Number != "" {
		return errors.New("Name already exists ")
	}

	err = repositories.SaveStudyPeriodName(studyModels.StudyPeriodName{
		Number:       uuid.NewV4().String(),
		Name:         studyPeriodNameRequest.Name,
		SchoolNumber: studyPeriodNameRequest.SchoolNumber,
		UserNumber:   currentUserNumber,
		Status:       constants.StatusActive,
	})

	if err != nil {
		return errors.New("Saving study period name error: " + err.Error())
	}

	return nil
}

func (studyPeriodService *studyPeriodService) FindAllStudyPeriodNames(schoolNumber string) ([]responses.StudyPeriodNameResponse, error) {

	var response = make([]responses.StudyPeriodNameResponse, 0)

	allNames, err := repositories.FindStudyPeriodName(schoolNumber)
	if err != nil {
		return response, errors.New("Fetching names error: " + err.Error())
	}

	for _, studyPeriodName := range allNames {
		response = append(response, createNameResponse(studyPeriodName))
	}

	return response, nil
}

func createNameResponse(studyPeriodName studyModels.StudyPeriodName) responses.StudyPeriodNameResponse {
	return responses.StudyPeriodNameResponse{
		DateAdded:             studyPeriodName.CreatedAt.Format(time.DateTime),
		StudyPeriodNameNumber: studyPeriodName.Number,
		Name:                  studyPeriodName.Name,
	}
}

func (studyPeriodService *studyPeriodService) UpdateStudyPeriodName(studyPeriodNameRequest requests.UpdateStudyPeriodNameRequest, currentUserNumber string) error {

	studyPeriodName, err := repositories.GetStudyPeriodNameByNumber(studyPeriodNameRequest.StudyPeriodNameNumber)
	if err != nil {
		return errors.New("Checking study period number error: " + err.Error())
	}

	checkName, err := repositories.GetStudyPeriodNameByName(studyPeriodNameRequest.SchoolNumber, studyPeriodNameRequest.Name)
	if err != nil {
		return errors.New("Checking name duplicate error " + err.Error())
	}

	if checkName.Number != "" && checkName.Number != studyPeriodName.Number {
		return errors.New("Name already exists ")
	}

	studyPeriodName.UserNumber = currentUserNumber
	studyPeriodName.Name = studyPeriodNameRequest.Name

	err = repositories.SaveStudyPeriodName(studyPeriodName)

	if err != nil {
		return errors.New("Updating study period name error: " + err.Error())
	}

	return nil
}

func (studyPeriodService *studyPeriodService) RemoveStudyPeriodName(studyPeriodNameNumber string) error {

	studyPeriodName, err := repositories.GetStudyPeriodNameByNumber(studyPeriodNameNumber)
	if err != nil {
		return errors.New("Checking study period number error: " + err.Error())
	}

	err = repositories.RemoveStudyPeriodName(studyPeriodName)
	if err != nil {
		return errors.New("Removing study period name error: " + err.Error())
	}

	return nil
}
