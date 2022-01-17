package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.fetchMovieById)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.fetchAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:genre_id", app.fetchMoviesByGenre)
	router.HandlerFunc(http.MethodGet, "/v1/genres", app.fetchAllGenres)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.insertMovie)
	router.HandlerFunc(http.MethodPut, "/v1/movie/:id", app.updateMovie)
	router.HandlerFunc(http.MethodDelete, "/v1/movie/:id", app.deleteMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies/search/:keyword", app.searchMovies)

	return app.enableCORS(router)
}
