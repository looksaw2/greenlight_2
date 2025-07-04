package api

import (
	"github.com/go-chi/chi/v5"
)

func (app *Application) NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/v1/health", app.ShowRealIPLog(app.HealthCheck))
	r.Post("/api/v1/create_movie", app.ShowRealIPLog(app.CreateMovieHandler))
	r.Get("/api/v1/show_movie/{id}", app.ShowRealIPLog(app.ShowMovieHandler))
	return r
}
