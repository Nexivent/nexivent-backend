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
	Eventos EventoSQL
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{
		Eventos: EventoSQL{DB: db},
	}
}
