package transport

import (
	"auth-service/internal/features/user/service"
)

type UserHandler struct {
	service service.UserService
}

func New(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}
