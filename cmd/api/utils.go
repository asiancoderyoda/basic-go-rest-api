package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) writeJSON(wr http.ResponseWriter, statusCode int, payload interface{}, wrap string) error {
	wrapper := make(map[string]interface{})
	wrapper[wrap] = payload
	js, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}
	wr.WriteHeader(statusCode)
	wr.Header().Set("Content-Type", "application/json")
	wr.Write(js)

	return nil
}

func (app *application) writeError(wr http.ResponseWriter, err error) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	httpError := errorResponse{
		Error: err.Error(),
	}

	app.writeJSON(wr, http.StatusUnprocessableEntity, httpError, "error")
}
