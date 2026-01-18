package transport

import "auth-service/internal/features/auth/service"

type AuthHandler struct {
	authService service.AuthService
}

func New(s service.AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}
