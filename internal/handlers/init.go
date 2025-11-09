package handlers

import (
	"github.com/govnocods/RedChat/internal/service"
)

type Handlers struct {
	UserService *service.UserService
}

func NewHandlers(userService *service.UserService) *Handlers {
	return &Handlers{
		UserService: userService,
	}
}
