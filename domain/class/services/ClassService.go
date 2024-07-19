package services

import (
	"errors"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/timetable/api/responses"
)

func FindSchoolClasses(schoolNumber, campusLevelNumber string) ([]responses.ClassResponse, error) {

	var classesResponse = make([]responses.ClassResponse, 0)

	classes, err := repositories.FindClassesBySchoolNumber(schoolNumber, campusLevelNumber)
	if err != nil {
		return classesResponse, errors.New("getting classes error: " + err.Error())
	}

	for _, class := range classes {
		classesResponse = append(classesResponse, CreateClassResponse(class))
	}

	return classesResponse, nil
}

func CreateClassResponse(class classModels.Class) responses.ClassResponse {
	return responses.ClassResponse{
		Number:       class.Number,
		Name:         class.Name,
		Abbreviation: class.Abbreviation,
	}
}
