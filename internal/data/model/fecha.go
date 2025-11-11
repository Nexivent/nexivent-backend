package model

import (
	"time"
)

type Fecha struct {
	ID          uint64
	FechaEvento time.Time

	EventoFechas []EventoFecha
}

func (Fecha) TableName() string { return "fecha" }
