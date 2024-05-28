package services

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"

type UserService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{repo: repo}
}
