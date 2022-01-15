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
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.insertMovie)

	return app.enableCORS(router)
}
