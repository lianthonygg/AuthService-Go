package service

import (
	"auth-service/internal/features/user/model"
	"auth-service/internal/features/user/store"
)

type UserService struct {
	store store.UserStore
}

func New(s store.UserStore) *UserService {
	return &UserService{store: s}
}

func (u *UserService) GetAllUsers() ([]*model.User, error) {
	users, err := u.store.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserService) GetUserById(id string) (*model.User, error) {
	user, err := u.store.GetById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) CreateUser(user *model.User) (*model.User, error) {
	user, err := u.store.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) UpdateUser(id string, user *model.User) (*model.User, error) {
	user, err := u.store.Update(id, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) RemoveUser(id string) error {
	err := u.store.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
