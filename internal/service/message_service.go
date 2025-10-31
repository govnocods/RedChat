package service

import "github.com/govnocods/RedChat/internal/repository"

type MessageService struct {
	repo repository.MessageRepositoryI
}

func NewMessageService(repo repository.MessageRepositoryI) *MessageService {
	return &MessageService{repo: repo}
}