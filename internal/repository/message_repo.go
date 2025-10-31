package repository

type MessageRepositoryI interface {

}

type messageRepository struct {

}

func NewMessageRepository() *messageRepository {
	return &messageRepository{}
}