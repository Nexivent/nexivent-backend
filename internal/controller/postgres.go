package controller

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBConfig struct {
	DNS             string
	MaxOpenConns    int
	MaxIdleConns    int
	MaxIdleTime     time.Duration // p.ej. 15 * time.Minute
	MaxConnLifetime time.Duration // p.ej. 2 * time.Hour
	PingTimeout     time.Duration // p.ej. 5 * time.Second
}

func OpenDB(conf DBConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", conf.DNS)
	if err != nil {
		return nil, err
	}

	// Defaults sensatos si vinieron en cero
	if conf.MaxOpenConns == 0 {
		conf.MaxOpenConns = 20
	}
	if conf.MaxIdleConns == 0 {
		conf.MaxIdleConns = 5
	}
	if conf.MaxIdleTime == 0 {
		conf.MaxIdleTime = 15 * time.Minute
	}
	if conf.MaxConnLifetime == 0 {
		conf.MaxConnLifetime = 2 * time.Hour
	}
	if conf.PingTimeout == 0 {
		conf.PingTimeout = 5 * time.Second
	}

	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetConnMaxIdleTime(conf.MaxIdleTime)
	db.SetConnMaxLifetime(conf.MaxConnLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), conf.PingTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
