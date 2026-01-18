package store

import (
	"database/sql"
	"errors"

	"auth-service/internal/features/user/model"
	"auth-service/internal/features/user/validate"
)

type UserStore interface {
	GetAll() ([]*model.User, error)
	GetById(id string) (*model.User, error)
	Create(user *validate.CreateUserRequest) (*model.User, error)
	Update(id string, user *model.User) (*model.User, error)
	Delete(id string) error
}

type store struct {
	db *sql.DB
}

func New(db *sql.DB) UserStore {
	return &store{db: db}
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

func (s *store) Create(user *validate.CreateUserRequest) (*model.User, error) {
	q := "INSERT INTO users (name, email, password) VALUES($1, $2, $3)"

	_, err := s.db.Exec(q, *user.Name, *user.Email, *user.Password)
	if err != nil {
		return nil, err
	}

	userModel := &model.User{
		Name:     *user.Name,
		Email:    *user.Email,
		Password: *user.Password,
	}

	return userModel, nil
}

func (s *store) Update(id string, user *model.User) (*model.User, error) {
	q := `UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4`

	_, err := s.db.Exec(q, user.Name, user.Email, user.Password, id)
	if err != nil {
		return nil, err
	}

	user.Id = id

	return user, nil
}

func (s *store) Delete(id string) error {
	q := "DELETE FROM users WHERE id = $1"

	_, err := s.db.Exec(q, id)
	if err != nil {
		return err
	}

	return nil
}
