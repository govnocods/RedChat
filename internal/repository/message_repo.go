package repository

import (
	"database/sql"
	"fmt"

	"github.com/govnocods/RedChat/models"
)

type MessageRepositoryI interface {
	Save(int, []byte) error
	GetMessage() ([]models.Message, error)
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

func (m *messageRepository) GetMessage() ([]models.Message, error) {
	query := `
		SELECT m.id, m.sender_id, m.content, m.created_at, u.username
		FROM (
			SELECT id, sender_id, content, created_at
			FROM messages
			ORDER BY created_at DESC
			LIMIT 50
		) m
		JOIN users u ON m.sender_id = u.id
		ORDER BY m.created_at ASC
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying messages: %w", err)
	}
	defer rows.Close()

	messages := make([]models.Message, 0, 50)

	for rows.Next() {
		var msg models.Message

		if err := rows.Scan(&msg.Id, &msg.SenderId, &msg.Content, &msg.CreatedAt, &msg.Username); err != nil {
			return nil, fmt.Errorf("error scanning message: %w", err)
		}

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	return messages, nil
}
