package api

import (
	"github.com/go-chi/chi/v5"
)

func (app *Application) NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.NotFound(app.notfoundResponse)
	r.MethodNotAllowed(app.methodNotAllowedResponse)
	r.Get("/api/v1/health", app.ShowRealIPLog(app.HealthCheck))
	r.Post("/api/v1/movie", app.ShowRealIPLog(app.CreateMovieHandler))
	r.Get("/api/v1/movie/{id}", app.ShowRealIPLog(app.ShowMovieHandler))
	r.Put("/api/v1/movie/{id}", app.UpdateMovieHandler)
	return r
}
