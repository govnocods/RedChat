package service

import (
	"github.com/govnocods/RedChat/internal/repository"
)

type MessageService struct {
	repo repository.MessageRepositoryI
}

func NewMessageService(repo repository.MessageRepositoryI) *MessageService {
	return &MessageService{repo: repo}
}

func (m *MessageService) SaveMessage(SenderId int, message []byte) error {
	return m.repo.Save(SenderId, message)
}
