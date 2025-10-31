package main

import (
	"net/http"

	"github.com/govnocods/RedChat/app"
	"github.com/govnocods/RedChat/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	database := database.NewDatabase()

	app := app.NewApp(database)

	go app.Hub.Run()

	http.ListenAndServe(":8080", app.Router)
}
