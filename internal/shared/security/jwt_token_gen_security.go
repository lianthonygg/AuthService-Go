package security

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"auth-service/internal/config"
	"auth-service/internal/features/user/model"

	"github.com/golang-jwt/jwt/v5"
)

type TokenGenerator interface {
	Generate(user *model.User) (string, error)
	GenerateRefreshToken() string
}

type JwtTokenGenerator struct {
	settings *config.Config
}

func NewGenerator(settings *config.Config) TokenGenerator {
	return &JwtTokenGenerator{settings: settings}
}

func (j *JwtTokenGenerator) Generate(user *model.User) (string, error) {
	now := time.Now().UTC()

	claims := jwt.MapClaims{
		"sub":   user.Id,
		"email": user.Email,
		"jti":   jwt.NewNumericDate(now).String(),
		"nbf":   now.Unix(),
		"iat":   now.Unix(),
		"iss":   j.settings.Issuer,
		"aud":   j.settings.Audience,
		"exp":   now.Add(time.Duration(j.settings.ExpirationHours) * time.Hour).Unix(),

		"role":        "User",
		"fullName":    user.Name,
		"isAvailable": true,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.settings.SecretKey))
}

func (j *JwtTokenGenerator) GenerateRefreshToken() string {
	b := make([]byte, 64)
	rand.Read(b)

	return hex.EncodeToString(b)
}
