package services

import (
	"errors"
	"fmt"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/api/responses"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
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

func FetchStreams(classNumber, campusNumber, curriculumNumber, campusCurriculumNumber, classLevelNumber string) ([]responses2.StreamInfo, error) {

	var streamResponses = make([]responses2.StreamInfo, 0)

	campusClass, err := ResolveCampusClass(classNumber, campusNumber, curriculumNumber, campusCurriculumNumber)

	if err != nil && classLevelNumber == "" {
		return streamResponses, err
	}

	if campusClass.Number == "" && classLevelNumber == "" {
		return streamResponses, errors.New("Selected class and curriculum are not on the selected campus ")
	} else {
		classLevelNumber = campusClass.Number
	}

	streams, err := repositories.FetchStreamByCampusClassNumber(campusClass.Number)
	if err != nil {
		return streamResponses, errors.New("Getting streams error " + err.Error())
	}

	for _, stream := range streams {
		streamResponses = append(streamResponses, responses2.StreamInfo{
			StreamNumber:      stream.Number,
			CampusClassNumber: stream.ClassLevelNumber,
			Name:              stream.Name,
		})
	}

	return streamResponses, nil
}
