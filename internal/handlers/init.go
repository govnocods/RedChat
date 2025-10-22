package handlers

import "github.com/govnocods/RedChat/internal/database"

type Handlers struct {
	DB *database.SQLDataBase
}

func NewHandlers(database *database.SQLDataBase) *Handlers {
	return &Handlers{
		DB: database,
	}
}