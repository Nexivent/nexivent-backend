package data

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Evento EventoModel
}

func NewModels(db *sqlx.DB) Models {
	return Models{
		Evento: EventoModel{DB: db},
	}
}
