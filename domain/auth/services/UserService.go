package services

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	authResponse "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/responses"
)

type UserService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func CreateUserResponse(user authModel.User) authResponse.UserResponse {
	return authResponse.UserResponse{
		Number:     user.Number,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		SecondName: user.SecondName,
		Gender:     user.Gender,
		Type:       user.Type,
		Status:     user.Status,
	}
}
