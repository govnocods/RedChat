package app

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/govnocods/RedChat/internal/handlers"
	"github.com/govnocods/RedChat/internal/middlewares"
	"github.com/govnocods/RedChat/internal/repository"
	"github.com/govnocods/RedChat/internal/service"
	"github.com/govnocods/RedChat/internal/websocket"
)

type App struct {
	DB          *sql.DB
	Router      *http.ServeMux
	Hub         *websocket.Hub
	Handlers    *handlers.Handlers
	Middlewares *middlewares.Middlewares
}

func NewApp(database *sql.DB) *App {
	userRepo := repository.NewUserRepository(database, context.Background())
	userService := service.NewUserService(userRepo)

	messageRepo := repository.NewMessageRepository(database)
	messageService := service.NewMessageService(messageRepo)

	app := &App{
		DB:          database,
		Router:      http.NewServeMux(),
		Hub:         websocket.NewHub(messageService),
		Handlers:    handlers.NewHandlers(userService),
		Middlewares: middlewares.NewMiddlewares(database),
	}

	app.routes()
	return app
}
