package service

import (
	"context"
	"errors"
	"time"

	authModel "auth-service/internal/features/auth/model"
	"auth-service/internal/features/auth/store"
	"auth-service/internal/features/user/model"
	"auth-service/internal/features/user/service"
	"auth-service/internal/features/user/validate"
	"auth-service/internal/shared/security"
)

type AuthService struct {
	userService  service.UserService
	authStore    store.AuthStore
	jwtGenerator security.TokenGenerator
	hasher       *security.PasswordHasher
}

func New(userService service.UserService, authStore store.AuthStore, jwtGenerator security.TokenGenerator, hasher *security.PasswordHasher) *AuthService {
	return &AuthService{userService: userService, authStore: authStore, jwtGenerator: jwtGenerator, hasher: hasher}
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

	refresh := a.jwtGenerator.GenerateRefreshToken()
	expires := time.Now().Add(7 * 24 * time.Hour)
	error := a.authStore.CreateRefreshToken(user.Id, refresh, expires)
	if error != nil {
		return nil, error
	}

	userResponse := authModel.UserResponse{
		Id:           user.Id,
		Name:         user.Name,
		Email:        user.Email,
		AccessToken:  token,
		RefreshToken: refresh,
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

	refresh := a.jwtGenerator.GenerateRefreshToken()
	expires := time.Now().Add(7 * 24 * time.Hour)
	error := a.authStore.CreateRefreshToken(created.Id, refresh, expires)
	if error != nil {
		return nil, error
	}

	userResponse := authModel.UserResponse{
		Id:           created.Id,
		Name:         created.Name,
		Email:        created.Email,
		AccessToken:  token,
		RefreshToken: refresh,
	}

	return &userResponse, nil
}

func (a *AuthService) Refresh(refreshTokenOld string) (*authModel.RefreshResponse, error) {
	ok, err := a.authStore.IsRevokeToken(refreshTokenOld)
	if ok {
		return nil, errors.New("Refresh Token Revoked")
	}

	refresh := a.jwtGenerator.GenerateRefreshToken()
	expires := time.Now().Add(7 * 24 * time.Hour)
	id, err := a.authStore.Rotate(refreshTokenOld, refresh, expires)
	if err != nil {
		return nil, err
	}

	user, err := a.GetUserById(id)
	if err != nil {
		return nil, err
	}

	token, err := a.jwtGenerator.Generate(user)
	if err != nil {
		return nil, err
	}

	refreshResponse := authModel.RefreshResponse{
		AccessToken:  token,
		RefreshToken: refresh,
	}

	return &refreshResponse, nil
}

func (a *AuthService) Logout(refreshToken string) (string, error) {
	err := a.authStore.Revoke(refreshToken)
	if err != nil {
		return "", err
	}

	return "Session closed successfully", nil
}
