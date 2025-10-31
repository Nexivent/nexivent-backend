package main

import (
	"log"
	"net/http"
	"os"
)

type application struct {
	logger *log.Logger
	conf   config
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	return app.logRequests(mux)
}
func main() {

	logger := log.New(os.Stdout, "", log.LstdFlags)
	conf := loadConfig(logger)

	app := &application{
		logger: logger,
		conf:   conf,
	}

	logger.Printf("listening on %s", conf.addr)
	if err := http.ListenAndServe(conf.addr, app.routes()); err != nil {
		logger.Fatal(err)
	}
}
