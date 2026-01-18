package store

import (
	"context"
	"database/sql"
	"errors"

	"auth-service/internal/features/user/model"
	"auth-service/internal/features/user/validate"
	"auth-service/internal/shared/security"

	"github.com/google/uuid"
)

type UserStore interface {
	GetAll() ([]*model.User, error)
	GetById(id string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Create(ctx context.Context, user *validate.CreateUserRequest) (*model.User, error)
	Update(id string, user *model.User) (*model.ResponseUserDTO, error)
	Delete(id string) error
}

type store struct {
	db     *sql.DB
	hasher *security.PasswordHasher
}

func New(db *sql.DB, hasher *security.PasswordHasher) UserStore {
	return &store{db: db, hasher: hasher}
}

func (s *store) GetAll() ([]*model.User, error) {
	q := "SELECT * FROM users"

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*model.User

	for rows.Next() {
		u := model.User{}
		if err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, &u)
	}

	if users == nil {
		users = []*model.User{}
	}

	return users, nil
}

func (s *store) GetById(id string) (*model.User, error) {
	q := "SELECT * FROM users WHERE id = $1"

	user := model.User{}
	err := s.db.QueryRow(q, id).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("cannot get user " + id + ": " + err.Error())
	}

	return &user, nil
}

func (s *store) GetByEmail(email string) (*model.User, error) {
	q := "SELECT * FROM users WHERE email = $1"

	user := model.User{}
	err := s.db.QueryRow(q, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("cannot get user " + email + ": " + err.Error())
	}

	return &user, nil
}

func (s *store) Create(ctx context.Context, user *validate.CreateUserRequest) (*model.User, error) {
	q := "INSERT INTO users (name, email, password) VALUES($1, $2, $3) RETURNING id, created_at"

	hashed, _ := s.hasher.Hash(*user.Password)

	var id uuid.UUID
	var createdAt string

	err := s.db.QueryRowContext(ctx, q, *user.Name, *user.Email, hashed).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	userModel := &model.User{
		Id:        id.String(),
		Name:      *user.Name,
		Email:     *user.Email,
		Password:  *user.Password,
		CreatedAt: createdAt,
	}

	return userModel, nil
}

func (s *store) Update(id string, user *model.User) (*model.ResponseUserDTO, error) {
	q := `UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4`

	_, err := s.db.Exec(q, user.Name, user.Email, user.Password, id)
	if err != nil {
		return nil, err
	}

	userModel := &model.ResponseUserDTO{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	return userModel, nil
}

func (s *store) Delete(id string) error {
	q := "DELETE FROM users WHERE id = $1"

	_, err := s.db.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}
