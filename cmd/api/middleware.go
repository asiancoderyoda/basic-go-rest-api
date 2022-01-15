package main

import "net/http"

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		wr.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(wr, r)
	})
}
