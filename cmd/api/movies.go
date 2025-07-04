package api

import (
	"fmt"
	"net/http"

	"github.com/looksaw/greenlight_2/internal/data"
	"github.com/looksaw/greenlight_2/internal/validator"
)

func (app *Application) CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Generes: input.Genres,
	}
	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.Models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) ShowMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readURLID(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	fmt.Fprintf(w, "show the details of movie %d\n", id)
}
