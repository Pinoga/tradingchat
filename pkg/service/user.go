package service

import (
	"tradingchat/pkg/model"
	"tradingchat/pkg/repo"
	"tradingchat/pkg/util"
)

type UserService interface {
	Register(username, password string) (string, error)
	FindByUsername(string) (*model.User, error)
	FindByID(string) (*model.User, error)
}

type userService struct {
	Repository repo.UserRepository
}

func NewUserService(repo repo.UserRepository) UserService {
	return &userService{
		Repository: repo,
	}
}

func (us *userService) Register(username, password string) (string, error) {
	hash, err := util.HashSaltPassword([]byte(password))
	if err != nil {
		return "", err
	}

	user := model.User{
		Username: username,
		Role:     "user",
		Hash:     hash,
	}

	id, err := us.Repository.Insert(&user)
	if err != nil {
		return "", err
	}
	return id, err
}

func (us *userService) FindByUsername(username string) (*model.User, error) {
	user, err := us.Repository.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (us *userService) FindByID(id string) (*model.User, error) {
	user, err := us.Repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, err
}
