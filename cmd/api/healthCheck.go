package api

import (
	"net"
	"net/http"
)

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	client_ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "client ip unknown, to avoid this , it must be valid url", http.StatusBadRequest)
		return
	}
	data := envelope{
		"status": "available",
		"system_info": map[string]any{
			"environment": app.Config.env,
			"version":     VERSION,
		},
	}
	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	app.Logger.Printf("client ip :%s\n", client_ip)
	app.Logger.Printf("visit is HealthCheck\n")
}
