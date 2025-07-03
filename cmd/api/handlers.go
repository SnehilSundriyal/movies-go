package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"watch-a-movie/internal/models"
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

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	// read json payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		log.Println(1)
		return
	}

	// validate user against database
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		log.Println(2)
		return
	}

	// check password
	valid, err := user.ValidatePassword(requestPayload.Password)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		log.Println(3)
		return
	}

	if !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		log.Println(4)
		return
	}

	// create a jwt user
	u := jwtUser{
		ID:        user.ID,        // Use actual user ID
		FirstName: user.FirstName, // Use actual user data
		LastName:  user.LastName,  // Use actual user data
	}

	// generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		log.Println(4)
		return
	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			// parse token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// get the user id from the token claims
			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			user, err := app.DB.GetUserByID(userID)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			u := jwtUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, app.auth.GetRefreshCookie(tokenPairs.RefreshToken))

			app.writeJSON(w, http.StatusOK, user)
		}
	}
}

func (app *application) displayMovie(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		ID int `json:"id"`
	}
	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	movie, err := app.DB.GetMovieByID(requestPayload.ID)
	if err != nil {
		app.errorJSON(w, errors.New("movie not found"), http.StatusNotFound)
		return
	}

	app.writeJSON(w, http.StatusOK, movie)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}

func (app *application) MovieCatalog(w http.ResponseWriter, r *http.Request) {
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

func (app *application) GetMovie(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie, err := app.DB.OneMovie(movieID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movie)
}

func (app *application) MovieForEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie, genres, err := app.DB.OneMovieForEdit(movieID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload = struct {
		Movie  *models.Movie   `json:"movie"`
		Genres []*models.Genre `json:"genres"`
	}{
		movie,
		genres,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) AllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.DB.AllGenres()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, genres)
}

func (app *application) InsertMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie

	err := app.readJSON(w, r, &movie)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if movie.IMDbID != "" {
		movie.IMDbID = extractIMDbIdFromLink(movie.IMDbID)

		if movie.IMDbID == "" {
			app.errorJSON(w, errors.New("invalid IMDb ID format"))
			return
		}
	}
	// try to get an image

	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()

	newID, err := app.DB.InsertMovie(movie)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// now handle genres
	err = app.DB.UpdateMovieGenres(newID, movie.GenresArray)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "movie updated",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func extractIMDbIdFromLink(imdbLink string) string {
	imdbLink = strings.TrimSpace(imdbLink)

	if strings.HasPrefix(imdbLink, "tt") && !strings.Contains(imdbLink, "/") {
		return imdbLink
	}

	re := regexp.MustCompile(`tt\d+`)
	match := re.FindString(imdbLink)

	return match
}
