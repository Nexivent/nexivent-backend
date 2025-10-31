package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type config struct {
	addr string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

// Carga variables desde el entorno; si existe .env, también lo lee.
func loadConfig(logger *log.Logger) config {
	_ = godotenv.Load() // no falla si .env no existe

	cfg := config{}
	cfg.addr = getEnv("ADDR", ":4000")

	cfg.db.dsn = getEnv("DB_DSN",
		"postgres://postgres:postgres@localhost:5432/nexivent?sslmode=disable")

	cfg.db.maxOpenConns = getEnvInt("DB_MAX_OPEN_CONNS", 25)
	cfg.db.maxIdleConns = getEnvInt("DB_MAX_IDLE_CONNS", 25)
	cfg.db.maxIdleTime = getEnvDuration("DB_MAX_IDLE_TIME", "15m")

	return cfg
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getEnvDuration(key, def string) time.Duration {
	v := getEnv(key, def)
	d, err := time.ParseDuration(v)
	if err != nil {
		dd, _ := time.ParseDuration(def)
		return dd
	}
	return d
}
