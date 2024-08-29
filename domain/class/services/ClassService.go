package services

import (
	"errors"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/repositories"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
)

func FindSchoolClasses(schoolNumber, campusLevelNumber, campusNumber, curriculumNumber string) ([]responses.ClassResponse, error) {

	var classesResponse = make([]responses.ClassResponse, 0)

	classes, err := repositories.FindClassesBySchoolNumber(schoolNumber, campusLevelNumber, campusNumber, curriculumNumber)
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

func ResolveCampusClass(classNumber, campusNumber, curriculumNumber string) (classModels.CampusClass, error) {
	campusClass, err := repositories.GetCampusClassByClassNumberAndCampusNumber(classNumber, campusNumber, "", "", curriculumNumber, "")
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
		return classModels.CampusClass{}, errors.New("Selected campus has this class on multiple curriculum, specify the curriculum ")
	}

	return campusClass[0], nil
}
