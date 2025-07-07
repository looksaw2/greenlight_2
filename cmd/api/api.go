package api

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/looksaw/greenlight_2/internal/data"
	"github.com/looksaw/greenlight_2/internal/jsonlog"
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
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type Application struct {
	Config Config
	Logger *jsonlog.Logger
	Models data.Models
}

func ApiInit() *Application {
	var cfg Config
	flag.IntVar(&cfg.port, "port", 8080, "This is the API server port")
	flag.StringVar(&cfg.env, "environment", "development", "This is about the environment")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL  DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PG max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PG max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PG max idle time")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum request per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rtae limiter maximun burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Parse()
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	app := Application{
		Config: cfg,
		Logger: logger,
	}
	return &app
}

func (app *Application) RunHTTP() error {
	r := app.NewRouter()
	db, err := openDB(app.Config)
	if err != nil {
		app.Logger.PrintInfo("Can't connect to the db", nil)
		app.Logger.PrintFatal(err, nil)
	}
	app.Logger.PrintInfo("connect to database ..............\n", nil)
	defer db.Close()
	app.Models = data.NewModel(db)
	runPortInfo := fmt.Sprintf("Server start to run on port : %d", app.Config.port)
	app.Logger.PrintInfo(runPortInfo, nil)
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	shutDownErr := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		app.Logger.PrintInfo("caught", map[string]string{
			"signal": s.String(),
		})
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		shutDownErr <- srv.Shutdown(ctx)
	}()
	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutDownErr
	if err != nil {
		return err
	}
	app.Logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})
	return nil
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
