package service

import (
	"tradingchat/pkg/store"
)

type UserService interface {
	FindByID(string) (*User, error)
	Register(username, password string) (string, error)
	FindByName(string) (*User, error)
}

type userService struct {
	Store store.Store
}

type User struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Role     string `json:"role" bson:"role"`
	Hash     string `json:"hash" bson:"hash"`
}

func NewUserService(store store.Store) UserService {
	return &userService{
		Store: store,
	}
}

func (us *userService) FindByID(id string) (*User, error) {
	var user User

	found, err := us.Store.FindOne("user", map[string]interface{}{"_id": id}, &user)
	if !found {
		return nil, err
	}
	return &user, err
}

func (us *userService) Register(username, password string) (string, error) {
	hash, err := HashSaltPassword([]byte(password))
	if err != nil {
		return "", err
	}

	user := User{
		Username: username,
		Role:     "user",
		Hash:     hash,
	}

	id, err := us.Store.InsertOne("user", &user)
	if err != nil {
		return "", err
	}
	return id, err
}

func (us *userService) FindByName(username string) (*User, error) {
	var user User
	found, err := us.Store.FindOne("user", map[string]interface{}{"username": username}, &user)
	if !found {
		return nil, err
	}
	return &user, err
}
