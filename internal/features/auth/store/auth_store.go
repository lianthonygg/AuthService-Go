package store

import (
	"database/sql"
	"time"
)

type AuthStore interface {
	IsRevokeToken(token string) (bool, error)
	CreateRefreshToken(userID, token string, expires time.Time) error
	Rotate(oldToken, newToken string, expiresAt time.Time) (string, error)
	Revoke(token string) error
}

type store struct {
	db *sql.DB
}

func NewAuthStore(db *sql.DB) AuthStore {
	return &store{db: db}
}

func (s *store) CreateRefreshToken(userID, token string, expires time.Time) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE refresh_tokens 
		SET revoked_at = now() 
		WHERE user_id = $1 AND revoked_at IS NULL
	`, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO refresh_tokens (user_id, token, expires_at, revoked_at)
		VALUES ($1, $2, $3, NULL)
	`, userID, token, expires)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *store) Rotate(oldToken, newToken string, expiresAt time.Time) (string, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE refresh_tokens
		SET revoked_at = now()
		WHERE token = $1
	`, oldToken)
	if err != nil {
		return "", err
	}

	var userID string
	err = tx.QueryRow(`
		INSERT INTO refresh_tokens (user_id, token, expires_at, revoked_at)
		SELECT user_id, $1, $2, NULL
		FROM refresh_tokens
		WHERE token = $3
		RETURNING user_id
	`, newToken, expiresAt, oldToken).Scan(&userID)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (s *store) IsRevokeToken(token string) (bool, error) {
	var revokedAt sql.NullTime
	err := s.db.QueryRow(`
		SELECT revoked_at FROM refresh_tokens WHERE token = $1
	`, token).Scan(&revokedAt)
	if err != nil {
		return false, err
	}

	return revokedAt.Valid, nil
}

func (s *store) Revoke(token string) error {
	_, err := s.db.Exec(`UPDATE refresh_tokens SET revoked_at = now() WHERE token=$1`, token)
	return err
}
