package models

import "time"

type Movie struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Poster         string    `json:"poster"`
	RuntimeHours   int       `json:"runtime"`
	RuntimeMinutes int       `json:"runtime_minutes"`
	IMDb           float32   `json:"imdb"`
	Release        int       `json:"release"`
	MPAA           string    `json:"mpaa"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
