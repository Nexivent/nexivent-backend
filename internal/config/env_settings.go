package config

import (
	"os"
	"strconv"

	"github.com/Loui27/nexivent-backend/logging"
	"github.com/Loui27/nexivent-backend/utils/env"
	"github.com/joho/godotenv"
)

type ConfigEnv struct {
	// LOGS
	EnableSqlLogs bool

	// SERVER
	MainPort      string
	EnableSwagger bool

	// DB
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDBName   string
	PostgresPsqlMode string
}

func NuevoConfigEnv(logger logging.Logger) *ConfigEnv {
	if ambiente := os.Getenv("NEXIVENT_POSTGRES_HOST"); ambiente == "local" || ambiente == "" {
		if envPath, err := env.FindEnvPath(); err != nil {
			logger.Panicln("Error finding .env file", err)
		} else if err := godotenv.Load(envPath); err != nil {
			logger.Panicln("Error loading .env file", err)
		}
	}
	enableSqlLogs, err := strconv.ParseBool(os.Getenv("ENABLE_SQL_LOGS"))
	if err != nil {
		enableSqlLogs = false
	}

	enableSwagger, err := strconv.ParseBool(os.Getenv("ENABLE_SWAGGER"))
	if err != nil {
		enableSwagger = false
	}

	mainPort := os.Getenv("MAIN_PORT")
	// Railway uses PORT environment variable
	if mainPort == "" {
		mainPort = os.Getenv("PORT")
	}
	// Default port if none is specified
	if mainPort == "" {
		mainPort = "8080"
	}

	PostgresHost := os.Getenv("NEXIVENT_POSTGRES_HOST")
	PostgresPort := os.Getenv("NEXIVENT_POSTGRES_PORT")
	PostgresUser := os.Getenv("NEXIVENT_POSTGRES_USER")
	PostgresPassword := os.Getenv("NEXIVENT_POSTGRES_PASSWORD")
	PostgresDBName := os.Getenv("NEXIVENT_POSTGRES_NAME")
	PostgresPsqlMode := os.Getenv("ASTRO_CAT_PSQL_SSL_MODE")

	return &ConfigEnv{
		EnableSqlLogs:    enableSqlLogs,
		MainPort:         mainPort,
		EnableSwagger:    enableSwagger,
		PostgresHost:     PostgresHost,
		PostgresPort:     PostgresPort,
		PostgresUser:     PostgresUser,
		PostgresPassword: PostgresPassword,
		PostgresDBName:   PostgresDBName,
		PostgresPsqlMode: PostgresPsqlMode,
	}
}
