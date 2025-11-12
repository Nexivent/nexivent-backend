package repository

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Repository struct {
	DB           *gorm.DB
	Eventos      EventoSQL
	EventoFechas EventoFecha
	Fechas       Fecha
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{
		DB:           db,
		Eventos:      EventoSQL{DB: db},
		EventoFechas: EventoFecha{DB: db},
		Fechas:       Fecha{DB: db},
	}
}

// WithTx crea una nueva instancia del Repository usando una transacción
func (r *Repository) WithTx(tx *gorm.DB) Repository {
	return Repository{
		DB:           tx,
		Eventos:      EventoSQL{DB: tx},
		EventoFechas: EventoFecha{DB: tx},
		Fechas:       Fecha{DB: tx},
	}
}
