package data

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBConfig struct {
	DNS          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

func OpenDB(dbConf DBConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbConf.DNS)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(dbConf.MaxOpenConns)
	db.SetMaxIdleConns(dbConf.MaxIdleConns)
	db.SetConnMaxIdleTime(dbConf.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
