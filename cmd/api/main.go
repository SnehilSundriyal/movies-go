package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"watch-a-movie/internal/repository"
	"watch-a-movie/internal/repository/dbrepo"
)

const port = 8080

type application struct {
	Domain string
	DSN    string
	DB     repository.DatabaseRepo
}

func main() {
	// set application config
	var app application

	// read from the command line
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=snehil password=hello dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.Parse()

	// connect to db
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	log.Println("Starting application on port:", port, "")

	// start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
