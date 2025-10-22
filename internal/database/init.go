package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/govnocods/RedChat/config"
)

type SQLDataBase struct {
	db *sql.DB
}

func NewDatabase() *SQLDataBase {
	var err error
	database, err := sql.Open("postgres", config.ConnStr)

	if err != nil {
		log.Fatal(err)
	}

	if err = database.Ping(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Successful connection to DataBase")
	}

	return &SQLDataBase{db: database}
}
