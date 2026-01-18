package service

import (
	"context"
	"errors"

	authModel "auth-service/internal/features/auth/model"
	"auth-service/internal/features/user/model"
	"auth-service/internal/features/user/service"
	"auth-service/internal/features/user/validate"
	"auth-service/internal/shared/security"
)

type AuthService struct {
	userService  service.UserService
	jwtGenerator security.TokenGenerator
	hasher       *security.PasswordHasher
}

func New(userService service.UserService, jwtGenerator security.TokenGenerator, hasher *security.PasswordHasher) *AuthService {
	return &AuthService{userService: userService, jwtGenerator: jwtGenerator, hasher: hasher}
}

func (a *AuthService) GetUserById(id string) (*model.User, error) {
	user, err := a.userService.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *AuthService) Login(email string, password string) (*authModel.UserResponse, error) {
	user, err := a.userService.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	ok := a.hasher.Verify(password, user.Password)
	if ok == false {
		return nil, errors.New("Invalid Credentials")
	}

	token, err := a.jwtGenerator.Generate(user)
	if err != nil {
		return nil, err
	}

	userResponse := authModel.UserResponse{
		Id:          user.Id,
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: token,
	}

	return &userResponse, nil
}

func (a *AuthService) Register(ctx context.Context, user *validate.CreateUserRequest) (*authModel.UserResponse, error) {
	created, err := a.userService.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := a.jwtGenerator.Generate(created)
	if err != nil {
		return nil, err
	}

	userResponse := authModel.UserResponse{
		Id:          created.Id,
		Name:        created.Name,
		Email:       created.Email,
		AccessToken: token,
	}

	return &userResponse, nil
}

func (a *AuthService) Refresh(refreshToken string) {
}
