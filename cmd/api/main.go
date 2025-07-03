package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
	"watch-a-movie/internal/repository"
	"watch-a-movie/internal/repository/dbrepo"
)

const port = 8080

type application struct {
	Domain       string
	DSN          string
	DB           repository.DatabaseRepo
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	APIKey       string
}

func main() {
	// set application config
	var app application

	// read from the command line
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=snehil password=hello dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")

	flag.StringVar(&app.JWTSecret, "jwt-secret", "very-secret", "signing secret for JWT")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer for JWT")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signing audience for JWT")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain for JWT")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain for JWT")
	flag.StringVar(&app.APIKey, "api-key", "3a159d69ea578f168c6504b7b56fe723", "api key")
	flag.Parse()

	// connect to db
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	app.auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "__Host-refresh-token",
		CookieDomain:  app.CookieDomain,
	}

	log.Println("Starting API on port 8080...")
	log.Println(`
  ______    ______       ______    ______   __
 /\  ___\  /\  __ \     /\  __ \  /\  == \ /\ \
 \ \ \__\\ \ \ \/\ \    \ \ \_\ \ \ \  __/ \ \ \
  \ \_____\ \ \_____\    \ \_\ \_\ \ \_\    \ \_\
   \/_____/  \/_____/     \/_/\/_/  \/_/     \/_/

`)

	// start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
