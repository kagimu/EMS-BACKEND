package services

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admissionLetters/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admissionLetters/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admissionLetters/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admissionLetters/repositories"
	uuid "github.com/satori/go.uuid"
)

type AdmissionLetterService struct{}

func NewAdmissionLetterService() AdmissionLetterService {
	return AdmissionLetterService{}
}

func (s AdmissionLetterService) Create(req requests.AdmissionLetterRequest) (responses.AdmissionLetterResponse, error) {
	letter := models.AdmissionLetter{
		Number:             uuid.NewV4().String(),
		SchoolNumber:      req.SchoolNumber,
		StudentNumber:     req.StudentNumber,
		CampusNumber:      req.CampusNumber,
		CampusClassNumber: req.CampusClassNumber,
		Phrase:            req.Phrase,
		PrincipalSignature: req.PrincipalSignature,
	}
	if err := repositories.SaveAdmissionLetter(&letter); err != nil {
		return responses.AdmissionLetterResponse{}, err
	}
	return s.toResponse(letter), nil
}

func (s AdmissionLetterService) Update(req requests.UpdateAdmissionLetterRequest) (responses.AdmissionLetterResponse, error) {
	letter, err := repositories.FindAdmissionLetterByNumber(req.Number)
	if err != nil {
		return responses.AdmissionLetterResponse{}, err
	}
	if letter.Number == "" {
		return responses.AdmissionLetterResponse{}, errors.New("admission letter not found")
	}
	letter.SchoolNumber = req.SchoolNumber
	letter.StudentNumber = req.StudentNumber
	letter.CampusNumber = req.CampusNumber
	letter.CampusClassNumber = req.CampusClassNumber
	letter.Phrase = req.Phrase
	letter.PrincipalSignature = req.PrincipalSignature
	if err := repositories.SaveAdmissionLetter(&letter); err != nil {
		return responses.AdmissionLetterResponse{}, err
	}
	return s.toResponse(letter), nil
}

func (s AdmissionLetterService) GetByNumber(number string) (responses.AdmissionLetterResponse, error) {
	letter, err := repositories.FindAdmissionLetterByNumber(number)
	if err != nil {
		return responses.AdmissionLetterResponse{}, err
	}
	if letter.Number == "" {
		return responses.AdmissionLetterResponse{}, errors.New("admission letter not found")
	}
	return s.toResponse(letter), nil
}

func (s AdmissionLetterService) GetByStudentNumber(studentNumber string) ([]responses.AdmissionLetterResponse, error) {
	letters, err := repositories.FindAdmissionLettersByStudentNumber(studentNumber)
	if err != nil {
		return nil, err
	}
	result := make([]responses.AdmissionLetterResponse, 0, len(letters))
	for _, l := range letters {
		result = append(result, s.toResponse(l))
	}
	return result, nil
}

func (s AdmissionLetterService) GetBySchoolNumber(schoolNumber string) ([]responses.AdmissionLetterResponse, error) {
	letters, err := repositories.FindAdmissionLettersBySchoolNumber(schoolNumber)
	if err != nil {
		return nil, err
	}
	result := make([]responses.AdmissionLetterResponse, 0, len(letters))
	for _, l := range letters {
		result = append(result, s.toResponse(l))
	}
	return result, nil
}

func (s AdmissionLetterService) Delete(number string) error {
	letter, err := repositories.FindAdmissionLetterByNumber(number)
	if err != nil {
		return err
	}
	if letter.Number == "" {
		return errors.New("admission letter not found")
	}
	return repositories.DeleteAdmissionLetter(&letter)
}

func (s AdmissionLetterService) toResponse(l models.AdmissionLetter) responses.AdmissionLetterResponse {
	return responses.AdmissionLetterResponse{
		Number:             l.Number,
		SchoolNumber:      l.SchoolNumber,
		StudentNumber:     l.StudentNumber,
		CampusNumber:      l.CampusNumber,
		CampusClassNumber: l.CampusClassNumber,
		Phrase:            l.Phrase,
		PrincipalSignature: l.PrincipalSignature,
		CreatedAt:         l.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		UpdatedAt:         l.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
	}
}
