package handlers

import (
	"github.com/govnocods/RedChat/internal/service"
)

type Handlers struct {
	UserService service.UserServiceI
}

func NewHandlers(userService service.UserServiceI) *Handlers {
	return &Handlers{
		UserService: userService,
	}
}
