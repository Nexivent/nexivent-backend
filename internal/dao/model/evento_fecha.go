package model

import (
	"time"
)

type EventoFecha struct {
	ID                  int64 `gorm:"column:evento_fecha_id;primaryKey;autoIncrement"`
	EventoID            int64
	FechaID             int64
	HoraInicio          time.Time
	Estado              int16
	UsuarioCreacion     *int64
	FechaCreacion       time.Time
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	Evento *Evento `gorm:"foreignKey:EventoID"`
	Fecha  *Fecha  `gorm:"foreignKey:FechaID"`

	Tickets []Ticket
}

func (EventoFecha) TableName() string { return "evento_fecha" }
