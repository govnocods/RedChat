package database

import (
	"database/sql"
	"fmt"

	"github.com/govnocods/RedChat/models"
)

func (db *SQLDataBase) CreateUser(user models.User) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`

	_, err := db.db.Exec(query, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("Error creating user: %w", err)
	}
	return nil
}

func (db *SQLDataBase) GetUser(username string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, username, password, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	row := db.db.QueryRow(query, username)

	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		return nil, err 
	}

	return &user, nil
}