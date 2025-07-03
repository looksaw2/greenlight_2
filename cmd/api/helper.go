package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type envelope map[string]any

func readURLID(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("Parse id is failed,Please input the correct id number")
	}
	if id <= 0 {
		return 0, errors.New("id must be positive")
	}
	return id, nil
}

func (app *Application) writeJSON(w http.ResponseWriter,
	status int,
	data any,
	headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}
