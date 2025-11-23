package config

import (
	"net/url"
	"os"
	"strconv"
	"strings"

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

	GoogleClientID string
}

func NuevoConfigEnv(logger logging.Logger) *ConfigEnv {
	if ambiente := os.Getenv("NEXIVENT_POSTGRES_HOST"); ambiente == "local" || ambiente == "" {
		if envPath, err := env.FindEnvPath(); err == nil {
			if err := godotenv.Load(envPath); err != nil {
				logger.Warnln("No se pudo cargar .env:", err)
			}
		} else if !os.IsNotExist(err) {
			logger.Warnln("Error buscando .env:", err)
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
	// Railway exposes the port through PORT
	if mainPort == "" {
		mainPort = os.Getenv("PORT")
	}
	if mainPort == "" {
		mainPort = "8098"
	}

	PostgresHost := os.Getenv("NEXIVENT_POSTGRES_HOST")
	PostgresPort := os.Getenv("NEXIVENT_POSTGRES_PORT")
	PostgresUser := os.Getenv("NEXIVENT_POSTGRES_USER")
	PostgresPassword := os.Getenv("NEXIVENT_POSTGRES_PASSWORD")
	PostgresDBName := os.Getenv("NEXIVENT_POSTGRES_NAME")
	PostgresPsqlMode := os.Getenv("ASTRO_CAT_PSQL_SSL_MODE")

	// Fallback to common PG_* variables (Railway/Postgres plugins)
	if PostgresHost == "" {
		PostgresHost = os.Getenv("PGHOST")
	}
	if PostgresPort == "" {
		PostgresPort = os.Getenv("PGPORT")
	}
	if PostgresUser == "" {
		PostgresUser = os.Getenv("PGUSER")
	}
	if PostgresPassword == "" {
		PostgresPassword = os.Getenv("PGPASSWORD")
	}
	if PostgresDBName == "" {
		PostgresDBName = os.Getenv("PGDATABASE")
	}
	if PostgresPsqlMode == "" {
		PostgresPsqlMode = os.Getenv("PGSSLMODE")
	}

	// Fallback to DATABASE_URL if provided
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		if parsed, err := url.Parse(dbURL); err == nil {
			if PostgresHost == "" {
				PostgresHost = parsed.Hostname()
			}
			if PostgresPort == "" {
				PostgresPort = parsed.Port()
			}
			if PostgresUser == "" && parsed.User != nil {
				PostgresUser = parsed.User.Username()
			}
			if PostgresPassword == "" && parsed.User != nil {
				if pwd, exists := parsed.User.Password(); exists {
					PostgresPassword = pwd
				}
			}
			if PostgresDBName == "" {
				PostgresDBName = strings.TrimPrefix(parsed.Path, "/")
			}
			if PostgresPsqlMode == "" {
				PostgresPsqlMode = parsed.Query().Get("sslmode")
			}
		} else {
			logger.Warnln("Could not parse DATABASE_URL:", err)
		}
	}

	// Default port if still unset
	if PostgresPort == "" {
		PostgresPort = "5432"
	}

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
		GoogleClientID:      os.Getenv("GOOGLE_CLIENT_ID"),
	}
}
