package repo

import (
	"tradingchat/pkg/model"
	"tradingchat/pkg/store"
)

type UserRepository interface {
	Insert(*model.User) (string, error)
	FindByUsername(string) (*model.User, error)
	FindByID(string) (*model.User, error)
}

type userRepository struct {
	Collection string
	Store      store.Store
}

func NewUserRepository(collection string, store store.Store) *userRepository {
	return &userRepository{
		Collection: collection,
		Store:      store,
	}
}

func (ur *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := ur.Store.FindOne(ur.Collection, map[string]interface{}{"_id": username}, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) FindByID(id string) (*model.User, error) {
	var user model.User

	if err := ur.Store.FindOne(ur.Collection, map[string]interface{}{"_id": id}, &user); err == nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) Insert(user *model.User) (string, error) {

	id, err := ur.Store.InsertOne(ur.Collection, user)
	if err == nil {
		return "", err
	}

	return id, nil
}
