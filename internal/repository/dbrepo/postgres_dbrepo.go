package dbrepo

import (
	"context"
	"database/sql"
	"strings"
	"time"
	"watch-a-movie/internal/models"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) AllMovies() ([]*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
    	SELECT 
    	    ID, TITLE, RUNTIME, IMDb, RELEASE, MPAA, DESCRIPTION, 
    	    COALESCE(poster, ''), CREATED_AT, UPDATED_AT 
    	FROM 
    	    MOVIES
    	ORDER BY
    	    title
    `

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*models.Movie

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.RuntimeHours,
			&movie.IMDb,
			&movie.Release,
			&movie.MPAA,
			&movie.Description,
			&movie.Poster,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if movie.Poster != "" {
			// Remove any whitespace/newlines and construct absolute URL
			cleanPoster := strings.TrimSpace(movie.Poster)
			movie.Poster = "http://localhost:8080/static/images/" + cleanPoster
		}

		movie.RuntimeMinutes = movie.RuntimeHours % 60
		movie.RuntimeHours = movie.RuntimeHours / 60

		movies = append(movies, &movie)
	}

	return movies, nil
}

func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT
			ID, EMAIL, FIRST_NAME, LAST_NAME, PASSWORD,
            CREATED_AT, UPDATED_AT 
		FROM
		    USERS
		WHERE
		    ID = $1
	`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT
			ID, EMAIL, FIRST_NAME, LAST_NAME, PASSWORD,
            CREATED_AT, UPDATED_AT 
		FROM
		    USERS
		WHERE
		    EMAIL = $1
	`

	var user models.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *PostgresDBRepo) GetMovieByID(id int) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT
			ID, TITLE, RUNTIME, IMDB, RELEASE, MPAA, DESCRIPTION, CREATED_AT, UPDATED_AT, POSTER, IMDB_ID
		FROM 
		    MOVIES
		WHERE 
		    ID = $1
`

	var movie models.Movie
	row := m.DB.QueryRowContext(ctx, query, id)

	var tempVar int

	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&tempVar,
		&movie.IMDb,
		&movie.Release,
		&movie.MPAA,
		&movie.Description,
		&movie.CreatedAt,
		&movie.UpdatedAt,
		&movie.Poster,
		&movie.IMDbID,
	)

	if err != nil {
		return nil, err
	}

	if movie.Poster != "" {
		cleanPoster := strings.TrimSpace(movie.Poster)
		movie.Poster = "http://localhost:8080/static/images/" + cleanPoster
	}

	movie.RuntimeMinutes = tempVar % 60
	movie.RuntimeHours = tempVar / 60

	return &movie, nil
}
