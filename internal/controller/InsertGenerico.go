package controller

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func InsertReturningID(cont context.Context, db ExecerQueryer, table string, cols []string, vals []any, idcol string, dest any) error {
	ph := make([]string, len(cols))
	for i := range cols {
		ph[i] = fmt.Sprintf("$%d", i+1)
	}
	q := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s",
		table, strings.Join(cols, ","), strings.Join(ph, ","), idcol)
	return db.QueryRowContext(cont, q, vals...).Scan(dest)

}

func UpdateByID(cont context.Context, db ExecerQueryer, table string, cols []string, vals []any, idcol string, id any) error {
	set := make([]string, len(cols))
	args := make([]any, 0, len(vals)+1)
	for i, c := range cols {
		set[i] = fmt.Sprintf("%s=$%d", c, i+1)
		args = append(args, vals[i])
	}
	args = append(args, id)
	q := fmt.Sprintf("UPDATE %s SET %s WHERE %s=$%d", table, strings.Join(set, ","), idcol, len(args))
	_, err := db.ExecContext(cont, q, args...)
	return err
}

type ExecerQueryer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
