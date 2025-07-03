package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (app *Application) CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is something about the createMovie")
}

func (app *Application) ShowMovieHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id must be a number", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "show the details of movie %d\n", id)
}
