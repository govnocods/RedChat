package handlers

import (
	"github.com/govnocods/RedChat/internal/service"
)

type Handlers struct {
	UserService service.UserServiceI
	MessageService service.MessageServiceI
}

func NewHandlers(userService service.UserServiceI, messageService service.MessageServiceI) *Handlers {
	return &Handlers{
		UserService: userService,
		MessageService: messageService, 
	}
}
