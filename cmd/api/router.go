package api

import (
	"github.com/go-chi/chi/v5"
)

func (app *Application) NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.NotFound(app.notfoundResponse)
	r.MethodNotAllowed(app.methodNotAllowedResponse)
	r.Get("/api/v1/health", app.ShowRealIPLog(app.HealthCheck))
	r.Get("/api/v1/movie", app.ShowRealIPLog(app.ListMovieHandler))
	r.Post("/api/v1/movie", app.ShowRealIPLog(app.CreateMovieHandler))
	r.Get("/api/v1/movie/{id}", app.ShowRealIPLog(app.ShowMovieHandler))
	r.Patch("/api/v1/movie/{id}", app.ShowRealIPLog(app.UpdateMovieHandler))
	r.Delete("/api/v1/movie/{id}", app.ShowRealIPLog(app.DeleteMovieHandler))
	return r
}
