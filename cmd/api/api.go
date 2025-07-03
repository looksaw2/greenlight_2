package api

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const VERSION = "1.0.0"

type Config struct {
	port int
	env  string
}

type Application struct {
	Config Config
	Logger *log.Logger
}

func ApiInit() *Application {
	var cfg Config
	flag.IntVar(&cfg.port, "port", 8080, "This is the API server port")
	flag.StringVar(&cfg.env, "environment", "development", "This is about the environment")
	flag.Parse()
	logger := log.New(os.Stdout, "[GreenlightLog] :", log.Ldate|log.Ltime)
	app := Application{
		Config: cfg,
		Logger: logger,
	}
	return &app
}

func (app *Application) RunHTTP() {
	r := app.NewRouter()
	app.Logger.Printf("Server start to run on port :%d", app.Config.port)
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err := srv.ListenAndServe()
	log.Fatal(err)
}
