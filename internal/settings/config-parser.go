package settings

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Configura de acuerdo a las opciones pasadas por comando o leídas desde el archivo .env
func ParseFlagEnv(logger *slog.Logger, cfg *Config) {
	var myEnv map[string]string
	myEnv, err := godotenv.Read()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.DB.URL, "db-url", myEnv["DATABASE_URL"], "PostgreSQL DSN")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()
}
