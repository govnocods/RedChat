package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var ConnStr string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env")
	}

	ConnStr = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)
}
