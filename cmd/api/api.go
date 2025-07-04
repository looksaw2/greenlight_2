package api

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const VERSION = "1.0.0"

type Config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type Application struct {
	Config Config
	Logger *log.Logger
}

func ApiInit() *Application {
	var cfg Config
	flag.IntVar(&cfg.port, "port", 8080, "This is the API server port")
	flag.StringVar(&cfg.env, "environment", "development", "This is about the environment")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL  DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PG max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PG max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PG max idle time")
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
	db, err := openDB(app.Config)
	if err != nil {
		app.Logger.Fatal(err)
	}
	app.Logger.Printf("connect to database ..............\n")
	defer db.Close()
	app.Logger.Printf("Server start to run on port :%d", app.Config.port)
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func openDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
