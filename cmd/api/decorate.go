package api

import (
	"net"
	"net/http"
)

func (app *Application) ShowRealIPLog(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//从请求中读取RealIP
		client_ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "httpAddr must be valid", http.StatusBadRequest)
			app.Logger.Printf("Parse realIP err %v", err)
			return
		}
		app.Logger.Printf("receive real ip is %s", client_ip)
		next(w, r)
		app.Logger.Printf("leave this middleware")
	}
}
