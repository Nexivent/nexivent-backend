package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Nexivent/nexivent-backend/internal/routing"
	"github.com/Nexivent/nexivent-backend/internal/settings"
)

func main() {
	var cfg settings.Config

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	settings.ParseFlagEnv(logger, &cfg)
	
	db, err := settings.OpenDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("database connection pool established")

	app := settings.Application{
		Config: cfg,
		Logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      routing.Routes(&app),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.Env)
	err = srv.ListenAndServe()
	logger.Error(err.Error())

	os.Exit(1)
}
