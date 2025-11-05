package model

import (
	"time"
)

type TipoDeTicket struct {
	ID                  int64 `gorm:"column:tipo_de_ticket_id;primaryKey;autoIncrement"`
	EventoID            int64
	Nombre              string
	FechaIni            time.Time
	FechaFin            time.Time
	Estado              int16
	UsuarioCreacion     int64
	FechaCreacion       time.Time
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	Evento *Evento `gorm:"foreignKey:EventoID"`
}

func (TipoDeTicket) TableName() string { return "tipo_de_ticket" }
