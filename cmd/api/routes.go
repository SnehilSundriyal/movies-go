package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"path/filepath"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.Home)
	mux.Get("/movies", app.AllMovies)
	mux.Post("/authenticate", app.authenticate)

	// Serve static files
	staticPath := filepath.Join("static")
	fileServer := http.FileServer(http.Dir(staticPath))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
