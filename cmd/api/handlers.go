package main

import (
	"log"
	"net/http"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payLoad = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Movies Go up and running",
		Version: "1.0.0",
	}

	err := app.writeJSON(w, http.StatusOK, payLoad)
	if err != nil {
		log.Println(err)
	}
}

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.AllMovies()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, movies)
	if err != nil {
		log.Println(err)
	}

}
