package config

import (
	"os"
	"strconv"

	"github.com/Nexivent/nexivent-backend/logging"
	"github.com/Nexivent/nexivent-backend/utils/env"
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

	// AWS S3
	AwsRegion           string
	AwsS3Bucket         string
	AwsS3Prefix         string
	AwsS3UploadDuration int64

	// Mail
	Host     string
	Port     int
	Username string
	Password string
	Sender   string

	// Factiliza
	FactilizaToken string `env:"FACTILIZA_TOKEN"`
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
		enableSwagger = true
	}

	mainPort := os.Getenv("MAIN_PORT")
	// Railway uses PORT environment variable
	if mainPort == "" {
		mainPort = os.Getenv("MAIN_PORT")
	}
	// Default port if none is specified
	if mainPort == "" {
		mainPort = "8098"
	}

	PostgresHost := os.Getenv("NEXIVENT_POSTGRES_HOST")
	PostgresPort := os.Getenv("NEXIVENT_POSTGRES_PORT")
	PostgresUser := os.Getenv("NEXIVENT_POSTGRES_USER")
	PostgresPassword := os.Getenv("NEXIVENT_POSTGRES_PASSWORD")
	PostgresDBName := os.Getenv("NEXIVENT_POSTGRES_NAME")
	PostgresPsqlMode := os.Getenv("ASTRO_CAT_PSQL_SSL_MODE")

	// AWS S3 config
	awsRegion := os.Getenv("AWS_REGION")
	awsBucket := os.Getenv("AWS_S3_BUCKET")
	awsPrefix := os.Getenv("AWS_S3_PREFIX")
	awsDurationStr := os.Getenv("AWS_S3_UPLOAD_EXPIRATION_SECONDS")
	var awsDuration int64 = 900 // default 15 minutes
	if awsDurationStr != "" {
		if v, err := strconv.ParseInt(awsDurationStr, 10, 64); err == nil && v > 0 {
			awsDuration = v
		}
	}

	// Mail config
	host := os.Getenv("MAIL_HOST")
	portStr := os.Getenv("MAIL_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 587 // default port
	}
	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	sender := os.Getenv("MAIL_SENDER")

	// Factiliza API token
	factilizaToken := os.Getenv("FACTILIZA_TOKEN")

	return &ConfigEnv{
		EnableSqlLogs:       enableSqlLogs,
		MainPort:            mainPort,
		EnableSwagger:       enableSwagger,
		PostgresHost:        PostgresHost,
		PostgresPort:        PostgresPort,
		PostgresUser:        PostgresUser,
		PostgresPassword:    PostgresPassword,
		PostgresDBName:      PostgresDBName,
		PostgresPsqlMode:    PostgresPsqlMode,
		AwsRegion:           awsRegion,
		AwsS3Bucket:         awsBucket,
		AwsS3Prefix:         awsPrefix,
		AwsS3UploadDuration: awsDuration,
		Host:                host,
		Port:                port,
		Username:            username,
		Password:            password,
		Sender:              sender,
		FactilizaToken:      factilizaToken,
	}
}
