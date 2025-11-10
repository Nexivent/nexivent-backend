package data

import (
	"time"

	"github.com/google/uuid"
)

type Fecha struct {
	ID          uuid.UUID
	FechaEvento time.Time

	EventoFechas []EventoFecha
}

func (Fecha) TableName() string { return "fecha" }
