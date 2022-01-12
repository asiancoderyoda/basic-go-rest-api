package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) statusHandler(rw http.ResponseWriter, r *http.Request) {
	currentStatus := ServerStatus{
		Version:     version,
		Status:      "OK",
		Environment: app.config.env,
	}

	js, err := json.MarshalIndent(currentStatus, "", "\t")
	if err != nil {
		app.logger.Println(err)
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(js)
}
