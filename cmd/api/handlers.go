package main

import (
	"encoding/json"
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

	out, err := json.Marshal(payLoad)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)
	if err != nil {
		log.Println(err)
	}

	w.Header()
}

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.AllMovies()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")

	out, err := json.Marshal(movies)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)
	if err != nil {
		log.Println(err)
	}

	w.Header()
}
