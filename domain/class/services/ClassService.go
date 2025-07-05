package services

import (
	"errors"
	"fmt"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/api/requests"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/api/responses"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"
	"time"
)

var admissionRepository = repositories4.NewAdmissionRepository()

func FindSchoolClasses(schoolNumber, campusLevelNumber, campusNumber, curriculumNumber, campusCurriculumNumber string) ([]responses.ClassResponse, error) {

	var classesResponse = make([]responses.ClassResponse, 0)

	classes, err := repositories.FindClassesBySchoolNumber(schoolNumber, campusLevelNumber, campusNumber, curriculumNumber, campusCurriculumNumber)
	if err != nil {
		return classesResponse, errors.New("getting classes error: " + err.Error())
	}

	for _, class := range classes {
		classesResponse = append(classesResponse, CreateClassResponse(class))
	}

	return classesResponse, nil
}

func CreateClassResponse(class classModels.Class) responses.ClassResponse {

	level, _ := repositories2.GetLevelByNumber(class.LevelNumber)

	return responses.ClassResponse{
		Number:       class.Number,
		Name:         class.Name,
		Abbreviation: class.Abbreviation,
		LevelNumber:  level.Number,
		LevelName:    level.Name,
		Code:         class.Code,
	}
}

func ResolveCampusClass(classNumber, campusNumber, curriculumNumber, campusCurriculumNumber string) (classModels.CampusClass, error) {
	campusClass, err := repositories.GetCampusClassByClassNumberAndCampusNumber(classNumber, campusNumber, curriculumNumber, campusCurriculumNumber)
	if err != nil {
		return classModels.CampusClass{}, errors.New("Checking class on campus error " + err.Error())
	}

	campus, err := repositories3.GetCampusByNumber(campusNumber)
	if err != nil {
		return classModels.CampusClass{}, errors.New("Checking campus error " + err.Error())
	}

	if campus.Number == "" {
		return classModels.CampusClass{}, errors.New("Invalid campus number ")
	}

	size := len(campusClass)

	if size == 0 {
		return classModels.CampusClass{}, errors.New("This class is not attached to the selected campus '" + campus.Name + "'")
	}

	if size > 1 {
		return classModels.CampusClass{}, errors.New(fmt.Sprintf("Selected campus has this class on multiple curriculum (%d), specify the curriculum %s %s %s", size, classNumber, campusNumber, curriculumNumber))
	}

	return campusClass[0], nil
}

func FetchStreams(classNumber, campusNumber, curriculumNumber, campusCurriculumNumber, classLevelNumber, campusClassNumber, status string, ctx *gin.Context) ([]responses2.StreamResponse, error) {

	var streamResponses = make([]responses2.StreamResponse, 0)
	var streams []classModels.Stream
	var err error
	if campusClassNumber != "" {
		streams, err = repositories.FetchStreamByCampusClassNumber(campusClassNumber, status, ctx)
		if err != nil {
			return streamResponses, errors.New("Getting streams error " + err.Error())
		}
	} else {

		campusClass, err := ResolveCampusClass(classNumber, campusNumber, curriculumNumber, campusCurriculumNumber)

		if err != nil && classLevelNumber == "" {
			return streamResponses, err
		}

		if campusClass.Number == "" && classLevelNumber == "" {
			return streamResponses, errors.New("Selected class and curriculum are not on the selected campus ")
		} else {
			classLevelNumber = campusClass.Number
		}

		streams, err = repositories.FetchStreamByCampusClassNumber(campusClass.Number, status, ctx)
		if err != nil {
			return streamResponses, errors.New("Getting streams error " + err.Error())
		}

	}
	for _, stream := range streams {
		streamResponses = append(streamResponses, CreateStreamResponse(stream))
	}

	return streamResponses, nil
}

func CreateStreamResponse(stream classModels.Stream) responses2.StreamResponse {
	return responses2.StreamResponse{
		StreamNumber:      stream.Number,
		CampusClassNumber: stream.CampusClassNumber,
		Name:              stream.Name,
		Status:            stream.Status,
	}
}

func FetchCampusClasses(ctx *gin.Context) ([]classModels.CampusClass, error) {

	campusClasses, err := repositories.FetchCampusClasses(ctx)
	if err != nil {
		return campusClasses, errors.New("Getting campus classes error " + err.Error())
	}

	return campusClasses, nil
}

func FetchCampusClassesByCampusNumber(campusNumber string) ([]classModels.CampusClass, error) {

	campusClasses, err := repositories.FindCampusClassesByCampusNumber(campusNumber)
	if err != nil {
		return campusClasses, errors.New("Getting campus classes error " + err.Error())
	}

	return campusClasses, nil
}

func UpdateStream(request requests.UpdateStreamRequest) error {
	stream, err := repositories.GetStreamByNumber(request.StreamNumber)
	if err != nil {
		return errors.New("Getting stream error " + err.Error())
	}

	if stream.Number == "" {
		return errors.New("Invalid Stream number ")
	}

	byName, err := repositories.GetStreamByNameAndCampusClassNumber(request.Name, stream.CampusClassNumber)
	if err != nil {
		return errors.New("Checking stream name error " + err.Error())
	}

	if byName.Number != "" && byName.Number != request.StreamNumber {
		return errors.New("Stream name already in use ")
	}

	status := stream.Status
	if status == "" {
		status = constants.StatusActive
	}

	if request.Status != "" {
		status = request.Status
	}

	stream.CampusClassNumber = request.CampusClassNumber
	stream.Name = request.Name
	stream.Status = status

	err = repositories.SaveStream(&stream)
	if err != nil {
		return errors.New("Updating stream error " + err.Error())
	}

	return nil
}

func DeleteStream(streamNumber string) error {
	stream, err := repositories.GetStreamByNumber(streamNumber)
	if err != nil {
		return errors.New("Getting stream error " + err.Error())
	}

	if stream.Number == "" {
		return errors.New("Invalid Stream number ")
	}

	admissions, err := admissionRepository.FindAdmissionsByStreamNumber(stream.Number)
	if err != nil {
		return errors.New("Checking admissions error " + err.Error())
	}

	admissionsCount := len(admissions)

	if admissionsCount > 0 {
		return errors.New(fmt.Sprintf("This stream has %d admissions, you can not delete it.", admissionsCount))
	}

	err = repositories.DeleteStreamByNumber(stream.Number)
	if err != nil {
		return errors.New("Deleting stream error " + err.Error())
	}

	return nil
}

func CreateStreams(streamRequest requests.StreamRequest, userNumber string) error {

	streams := strings.Split(streamRequest.Name, ",")

	//add streams to the class
	if streams != nil {
		for _, stream := range streams {

			streamObj := classModels.Stream{
				Number:            uuid.NewV4().String(),
				Name:              stream,
				CampusClassNumber: streamRequest.CampusClassNumber,
				UserNumber:        userNumber,
				Status:            constants.StatusActive,
			}

			byName, err := repositories.GetStreamByNameAndCampusClassNumber(stream, streamRequest.CampusClassNumber)
			if err != nil {
				return errors.New("Checking duplicate stream error " + err.Error())
			}

			if byName.Number == "" {
				err := repositories.SaveStream(&streamObj)
				if err != nil {
					return errors.New("Creating stream error " + err.Error())
				}
			}

		}
	}

	return nil
}

func StreamDetails(streamNumber string) (responses2.StreamResponse, error) {

	stream, err := repositories.GetStreamByNumber(streamNumber)
	if err != nil {
		return responses2.StreamResponse{}, errors.New("Getting stream error " + err.Error())
	}

	if stream.Number == "" {
		return responses2.StreamResponse{}, errors.New("Invalid Stream number ")
	}

	return CreateStreamResponse(stream), nil
}

func EditClassDetails(info requests.ClassInfo, userNumber string) error {
	class, err := repositories.GetClassByNumber(info.ClassNumber)
	if err != nil {
		return errors.New("Getting class error " + err.Error())
	}

	if class.Number == "" {
		return errors.New("Invalid Class number ")
	}

	codeNo, err := strconv.ParseInt(info.Code, 10, 64)
	if err != nil {
		return errors.New("Invalid Code: " + err.Error())
	}

	classByCode, err := repositories.GetClassByCode(codeNo, info.EntityNumber)
	if err != nil {
		return errors.New("Getting class by code error " + err.Error())
	}

	if classByCode.Number != "" && classByCode.Number != class.Number {
		return errors.New("Code already in use ")
	}

	classByName, err := repositories.GetClassByName(info.Name, info.EntityNumber)
	if err != nil {
		return errors.New("Getting class by name error " + err.Error())
	}

	if classByName.Number != "" && classByName.Number != class.Number {
		return errors.New("Name already in use ")
	}

	classByAbbreviation, err := repositories.GetClassByAbbreviation(info.Abbreviation, info.EntityNumber)
	if err != nil {
		return errors.New("Getting class by abbreviation error " + err.Error())
	}

	if classByAbbreviation.Number != "" && classByAbbreviation.Number != class.Number {
		return errors.New("Abbreviation already in use ")
	}

	return repositories.SaveClass(classModels.Class{
		Number:       class.Number,
		Name:         info.Name,
		Code:         codeNo,
		EntityNumber: classByCode.Number,
		Abbreviation: classByAbbreviation.Number,
		CreatedAt:    class.CreatedAt,
		UpdatedAt:    time.Now(),
		UserNumber:   userNumber,
		LevelNumber:  class.LevelNumber,
	})
}
