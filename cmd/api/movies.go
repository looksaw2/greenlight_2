package api

import (
	"fmt"
	"net/http"
)

func (app *Application) CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is something about the createMovie")
}

func (app *Application) ShowMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readURLID(r)
	if err != nil {
		http.Error(w, "id must be a number", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "show the details of movie %d\n", id)
}
