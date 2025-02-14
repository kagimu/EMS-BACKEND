package services

import (
	"errors"
	"fmt"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/api/requests"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/api/responses"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
	"github.com/gin-gonic/gin"
)

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

func FetchStreams(classNumber, campusNumber, curriculumNumber, campusCurriculumNumber, classLevelNumber, campusClassNumber string, ctx *gin.Context) ([]responses2.StreamResponse, error) {

	var streamResponses = make([]responses2.StreamResponse, 0)
	var streams []classModels.Stream
	var err error
	if campusClassNumber != "" {
		streams, err = repositories.FetchStreamByCampusClassNumber(campusClassNumber, ctx)
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

		streams, err = repositories.FetchStreamByCampusClassNumber(campusClass.Number, ctx)
		if err != nil {
			return streamResponses, errors.New("Getting streams error " + err.Error())
		}

	}
	for _, stream := range streams {
		streamResponses = append(streamResponses, CreateStreamResponse(stream))
	}

	return streamResponses, nil
}

func FetchStreamsByCampusClassNumber(campusClassNumber string, ctx *gin.Context) ([]responses2.StreamResponse, error) {
	var streamResponses = make([]responses2.StreamResponse, 0)

	streams, err := repositories.FetchStreamByCampusClassNumber(ctx.Query("campus_class_number"), ctx)
	if err != nil {
		return streamResponses, errors.New("Getting streams error " + err.Error())
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

	stream.CampusClassNumber = request.CampusClassNumber
	stream.Name = request.Name

	err = repositories.SaveStream(&stream)
	if err != nil {
		return errors.New("Updating stream error " + err.Error())
	}

	return nil
}
