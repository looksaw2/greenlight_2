package api

import (
	"errors"
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
	movie, err := app.Models.Movies.Get(int64(id))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) UpdateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readURLID(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}
	movie, err := app.Models.Movies.Get(int64(id))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:runtime"`
		Genres  []string     `json:"genres"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Generes = input.Genres

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.Models.Movies.Update(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
func (app *Application) DeleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := readURLID(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}
	err = app.Models.Movies.Delete(int64(id))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"msg": "The movie has been success delete"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, nil)
	}
}
