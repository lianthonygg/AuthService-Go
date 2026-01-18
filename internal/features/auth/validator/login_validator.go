package validator

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

var (
	ErrMissingEmail    = errors.New("email is required")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrMissingPassword = errors.New("password is required")
	ErrInvalidPassword = errors.New("password does not meet complexity requirements")
)

func (r *LoginRequest) Validate() error {
	var errs []error

	if r.Email == nil || strings.TrimSpace(*r.Email) == "" {
		errs = append(errs, ErrMissingEmail)
	} else if !strings.Contains(*r.Email, "@") {
		errs = append(errs, ErrInvalidEmail)
	}

	if r.Password == nil || strings.TrimSpace(*r.Password) == "" {
		errs = append(errs, ErrMissingPassword)
	} else if len(*r.Password) < 8 {
		errs = append(errs, ErrInvalidPassword)
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}

func DecodeAndValidateJSON[T any](r *http.Request, dst *T) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("content-type must be application/json")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if !json.Valid(body) {
		return errors.New("invalid JSON format")
	}

	if err := json.Unmarshal(body, dst); err != nil {
		return err
	}

	if v, ok := any(dst).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}
