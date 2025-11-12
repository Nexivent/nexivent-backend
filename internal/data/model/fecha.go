package model

import (
	"time"
)

type Fecha struct {
	ID          uint64    `gorm:"column:fecha_id;primaryKey" json:"id"`
	FechaEvento time.Time `gorm:"column:fecha_evento" json:"fechaEvento"`

	EventoFechas []EventoFecha `json:"-"`
}

func (Fecha) TableName() string { return "fecha" }
