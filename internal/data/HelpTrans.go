package data

import (
	"context"
	"database/sql"
)

func WithTx(cont context.Context, db *sql.DB, fn func(context.Context, *sql.Tx) error) error {
	tx, err := db.BeginTx(cont, nil)
	if err != nil {
		return err
	}
	if err := fn(cont, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
