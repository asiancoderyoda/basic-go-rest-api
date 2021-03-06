package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) fetchMovieById(wr http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Println(errors.New("invalid movie id " + params.ByName("id")))
		app.writeError(wr, err)
		return
	}

	movie, err := app.models.DB.GetMovieById(id)
	if err != nil {
		app.logger.Fatalf("Error while quering data: %v", err)
		return
	}

	err = app.writeJSON(wr, http.StatusOK, movie, "movie")
	if err != nil {
		app.logger.Fatalf("Error while writing response: %v", err)
		return
	}
}

func (app *application) fetchAllMovies(wr http.ResponseWriter, r *http.Request) {
	movies, err := app.models.DB.GetAllMovie()
	if err != nil {
		app.logger.Fatalf("Error while quering data: %v", err)
		return
	}
	err = app.writeJSON(wr, http.StatusOK, movies, "movies")
	if err != nil {
		app.logger.Fatalf("Error while writing response: %v", err)
		return
	}
}

func (app *application) fetchMoviesByGenre(wr http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	genre_id, err := strconv.Atoi(params.ByName("genre_id"))
	if err != nil {
		app.logger.Println(errors.New("invalid movie id " + params.ByName("genre_id")))
		app.writeError(wr, err)
		return
	}
	movies, err := app.models.DB.GetAllMovie(genre_id)
	if err != nil {
		app.logger.Fatalf("Error while quering data: %v", err)
		return
	}
	err = app.writeJSON(wr, http.StatusOK, movies, "movies")
	if err != nil {
		app.logger.Fatalf("Error while writing response: %v", err)
		return
	}
}

func (app application) fetchAllGenres(wr http.ResponseWriter, r *http.Request) {
	genres, err := app.models.DB.GetAllGenres()
	if err != nil {
		app.logger.Fatalf("Error while quering data: %v", err)
	}
	err = app.writeJSON(wr, http.StatusOK, genres, "genres")
	if err != nil {
		app.logger.Fatalf("Error while writing response: %v", err)
	}
}

func (app *application) insertMovie(wr http.ResponseWriter, r *http.Request) {
	// var movie Movie
	// err := json.NewDecoder(r.Body).Decode(&movie)
	// if err != nil {
	// 	app.logger.Println(errors.New("invalid movie data"))
	// 	app.writeError(wr, err)
	// 	return
	// }

	// err = app.models.DB.CreateMovie(movie)
	// if err != nil {
	// 	app.logger.Fatalf("Error while quering data: %v", err)
	// 	return
	// }

	// err = app.writeJSON(wr, http.StatusCreated, movie, "movie")
	// if err != nil {
	// 	app.logger.Fatalf("Error while writing response: %v", err)
	// 	return
	// }
}

func (app *application) updateMovie(wr http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteMovie(wr http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Println(errors.New("invalid movie id " + params.ByName("id")))
		app.writeError(wr, err)
		return
	}
	success, err := app.models.DB.DeleteMovie(id)
	if err != nil {
		app.logger.Fatalf("Error while deleting movie: %v", err)
		return
	}
	if !success {
		app.writeError(wr, errors.New("movie not found"))
		return
	}
	app.writeJSON(wr, http.StatusOK, nil, "")
}

func (app *application) searchMovies(wr http.ResponseWriter, r *http.Request) {

}
