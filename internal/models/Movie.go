package models

import "time"

type Movie struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Poster         string    `json:"poster"`
	RuntimeHours   int       `json:"runtime"`
	RuntimeMinutes int       `json:"runtime_minutes"`
	IMDb           float32   `json:"imdb"`
	IMDbID         string    `json:"imdbId"`
	Release        int       `json:"release"`
	MPAA           string    `json:"mpaa"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
	Genres         []*Genre  `json:"genres,omitempty"`
	GenresArray    []int     `json:"genres_array,omitempty"`
}

type Genre struct {
	ID        int       `json:"id"`
	Genre     string    `json:"genre"`
	Checked   bool      `json:"checked"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
