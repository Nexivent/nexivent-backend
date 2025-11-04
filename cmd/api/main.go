package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	//model "github.com/Loui27/nexivent-backend/internal/dao/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	logger *log.Logger
	conf   config
	db     *sql.DB
	//categorias model.CategoriaRepository
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
