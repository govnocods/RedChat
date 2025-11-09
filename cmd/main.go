package main

import (
	"net/http"
	"os"

	"github.com/govnocods/RedChat/app"
	"github.com/govnocods/RedChat/internal/database"
	"github.com/govnocods/RedChat/internal/logger"
	_ "github.com/lib/pq"
)

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logger.InitLogger(logger.ParseLevel(logLevel))

	database := database.NewDatabase()

	app := app.NewApp(database)

	go app.Hub.Run()

	logger.Info("Server listening", "address", ":8080")
	if err := http.ListenAndServe(":8080", app.Router); err != nil {
		logger.WithError(err).Error("Server failed to start")
		os.Exit(1)
	}
}
