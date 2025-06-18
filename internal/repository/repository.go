package repository

import (
	"database/sql"
	"watch-a-movie/internal/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	AllMovies() ([]*models.Movie, error)
}
