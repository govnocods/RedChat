package database

import (
	"database/sql"
	"os"

	"github.com/govnocods/RedChat/config"
	"github.com/govnocods/RedChat/internal/logger"
)

type SQLDataBase struct {
	db *sql.DB
}

func NewDatabase() *sql.DB {
	var err error
	database, err := sql.Open("postgres", config.ConnStr)

	if err != nil {
		logger.WithError(err).Error("Failed to open database connection")
		os.Exit(1)
	}

	if err = database.Ping(); err != nil {
		logger.WithError(err).Error("Failed to ping database")
		os.Exit(1)
	}

	logger.Info("Database connection established")

	return database
}
