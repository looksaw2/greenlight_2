package api

import "net/http"

func (app *Application) LogError(r *http.Request, err error) {
	app.Logger.Print(err)
}

func (app *Application) errorResponse(w http.ResponseWriter,
	r *http.Request,
	status int,
	msg any) {

	env := envelope{"msg": msg}
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.LogError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *Application) serverErrorResponse(w http.ResponseWriter,
	r *http.Request,
	err error) {
	app.LogError(r, err)
	message := "The server meet some problem so couldn't receive your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *Application) notfoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The Page is not found , please check your url"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := "This method is not allowed , please use the legal method"
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadGateway, err.Error())
}

func (app *Application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
