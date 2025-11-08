package repository

import (
	"database/sql"
	"fmt"
)

type MessageRepositoryI interface {
	Save(int, []byte) error
}

type messageRepository struct {
	DB *sql.DB
}

func NewMessageRepository(database *sql.DB) *messageRepository {
	return &messageRepository{
		DB: database,
	}
}

func (m *messageRepository) Save(senderID int, message []byte) error {
	query := `INSERT INTO messages (sender_id, content) VALUES ($1, $2)`

	_, err := m.DB.Exec(query, senderID, message)

	if err != nil {
		return fmt.Errorf("Error save message: %w", err)
	}

	return nil
}
