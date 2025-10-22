package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/govnocods/RedChat/models"
)

type UserRepositoryI interface {
	CreateUser(user models.User) error
	GetUser(username string) (*models.User, error)
}

type userRepository struct {
	DB  *sql.DB
	Ctx context.Context
}

func NewUserRepository(db *sql.DB, ctx context.Context) *userRepository {
	return &userRepository{
		DB:  db,
		Ctx: ctx,
	}
}

func (r *userRepository) CreateUser(user models.User) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`

	_, err := r.DB.Exec(query, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("Error creating user: %w", err)
	}
	return nil
}

func (r *userRepository) GetUser(username string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, username, password, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	row := r.DB.QueryRow(query, username)

	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
