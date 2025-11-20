package service

import (
	"github.com/govnocods/RedChat/internal/repository"
	"github.com/govnocods/RedChat/models"
)

type MessageServiceI interface {
	SaveMessage(SenderId int, message []byte) error
	GetMessage() ([]models.Message, error)
}

type MessageService struct {
	repo repository.MessageRepositoryI
}

func NewMessageService(repo repository.MessageRepositoryI) *MessageService {
	return &MessageService{repo: repo}
}

func (m *MessageService) SaveMessage(SenderId int, message []byte) error {
	return m.repo.Save(SenderId, message)
}

func (m *MessageService) GetMessage() ([]models.Message, error) {
	messages, err := m.repo.GetMessage()
	if err != nil {
		return nil, err
	}

	return messages, err
}
