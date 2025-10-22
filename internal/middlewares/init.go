package middlewares

import "database/sql"

type Middlewares struct {
	DB *sql.DB
}

func NewMiddlewares(database *sql.DB) *Middlewares {
	return &Middlewares{
		DB: database,
	}
}
