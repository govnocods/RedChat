package middlewares

import "github.com/govnocods/RedChat/internal/database"

type Middlewares struct {
	DB *database.SQLDataBase
}

func NewMiddlewares(database *database.SQLDataBase) *Middlewares {
	return &Middlewares{
		DB: database,
	}
}