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
    	    COALESCE(poster, ''), CREATED_AT, UPDATED_AT, IMDB_ID
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
			&movie.IMDbID,
		)
		if err != nil {
			return nil, err
		}
		if movie.Poster != "" {
			cleanPoster := strings.TrimSpace(movie.Poster)
			movie.Poster = "http://localhost:8080/static/images/" + cleanPoster
		}

		movie.RuntimeMinutes = movie.RuntimeHours % 60
		movie.RuntimeHours = movie.RuntimeHours / 60

		movies = append(movies, &movie)
	}

	return movies, nil
}

func (m *PostgresDBRepo) OneMovie(id int) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			id, title, release, runtime, imdb, mpaa, description,
			COALESCE(poster, ''), created_at, updated_at, imdb_id
		FROM 
		    MOVIES
		WHERE
		    ID = $1
    `

	row := m.DB.QueryRowContext(ctx, query, id)

	var movie models.Movie
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Release,
		&movie.RuntimeHours,
		&movie.IMDb,
		&movie.MPAA,
		&movie.Description,
		&movie.Poster,
		&movie.CreatedAt,
		&movie.UpdatedAt,
		&movie.IMDbID,
	)
	if err != nil {
		return nil, err
	}

	if movie.Poster != "" {
		cleanPoster := strings.TrimSpace(movie.Poster)
		movie.Poster = "http://localhost:8080/static/images/" + cleanPoster
	}

	movie.RuntimeMinutes = movie.RuntimeHours % 60
	movie.RuntimeHours = movie.RuntimeHours / 60

	// get genres, if any
	query = `
		SELECT 
		    g.id, g.genre
		FROM 
		    MOVIES_GENRES mg 
		LEFT JOIN 
			GENRES g 
		ON 
			(mg.genre_id = g.id)
		WHERE
		    mg.movie_id = $1
		ORDER BY
		    g.genre
`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	var genres []*models.Genre
	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, err
		}

		genres = append(genres, &g)
	}

	movie.Genres = genres

	return &movie, err
}

func (m *PostgresDBRepo) OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			id, title, release, runtime, imdb, mpaa, description,
			COALESCE(poster, ''), created_at, updated_at, imdb_id
		FROM 
		    MOVIES
		WHERE
		    ID = $1
    `

	row := m.DB.QueryRowContext(ctx, query, id)

	var movie models.Movie
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Release,
		&movie.RuntimeHours,
		&movie.IMDb,
		&movie.MPAA,
		&movie.Description,
		&movie.Poster,
		&movie.CreatedAt,
		&movie.UpdatedAt,
		&movie.IMDbID,
	)
	if err != nil {
		return nil, nil, err
	}

	if movie.Poster != "" {
		cleanPoster := strings.TrimSpace(movie.Poster)
		movie.Poster = "http://localhost:8080/static/images/" + cleanPoster
	}

	movie.RuntimeMinutes = movie.RuntimeHours % 60
	movie.RuntimeHours = movie.RuntimeHours / 60

	// get genres, if any
	query = `
		SELECT 
		    g.id, g.genre
		FROM 
		    MOVIES_GENRES mg 
		LEFT JOIN 
			GENRES g 
		ON 
			(mg.genre_id = g.id)
		WHERE
		    mg.movie_id = $1
		ORDER BY
		    g.genre
`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, nil, err
	}
	defer rows.Close()

	var genres []*models.Genre
	var genresArray []int
	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, nil, err
		}

		genres = append(genres, &g)
		genresArray = append(genresArray, g.ID)
	}

	movie.Genres = genres
	movie.GenresArray = genresArray

	var allGenres []*models.Genre
	query = `
		SELECT 
		    ID, GENRE
		FROM 
		    GENRES
		ORDER BY
		    GENRE
    `
	gRows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	defer gRows.Close()

	for gRows.Next() {
		var g models.Genre
		err := gRows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, nil, err
		}

		allGenres = append(allGenres, &g)
	}

	return &movie, allGenres, err

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

func (m *PostgresDBRepo) AllGenres() ([]*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
			ID, GENRE, CREATED_AT, UPDATED_AT
		FROM
		    GENRES
		ORDER BY
		    GENRE
`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []*models.Genre

	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		genres = append(genres, &g)
	}

	return genres, nil
}

func (m *PostgresDBRepo) InsertMovie(movie models.Movie) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		INSERT INTO
			MOVIES (title, description, release, runtime, mpaa, imdb,)
`
}
