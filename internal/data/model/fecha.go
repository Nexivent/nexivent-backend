package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type Fecha struct {
	ID          uint64    `gorm:"column:fecha_id;primaryKey;autoIncrement" json:"id"`
	FechaEvento time.Time `gorm:"column:fecha_evento;unique" json:"fechaEvento"`

	EventoFechas []EventoFecha `json:"-"`
}

func (Fecha) TableName() string { return "fecha" }

func ValidateFecha(v *validator.Validator, fecha *Fecha) {
	// Validar FechaEvento
	v.Check(!fecha.FechaEvento.IsZero(), "fechaEvento", "la fecha del evento es obligatoria")
}
