package app

import (
	"net/http"

	"github.com/govnocods/RedChat/internal/database"
	"github.com/govnocods/RedChat/internal/handlers"
	"github.com/govnocods/RedChat/internal/middlewares"
)

type App struct {
	DB      *database.SQLDataBase
	Router  *http.ServeMux
	Handlers *handlers.Handlers
	Middlewares *middlewares.Middlewares
}

func NewApp(database *database.SQLDataBase) *App {
	app := &App{
		DB: database,
		Router: http.NewServeMux(),
		Handlers: handlers.NewHandlers(database),
		Middlewares: middlewares.NewMiddlewares(database),
	}

	app.routes()
	return app
}
