package model

import (
	"time"
)

type TipoDeTicket struct {
	ID                  int64  `gorm:"column:tipo_de_ticket_id;primaryKey;autoIncrement"`
	EventoID            int64  `gorm:"uniqueIndex:uq_tipo_ticket_nombre"`
	Nombre              string `gorm:"uniqueIndex:uq_tipo_ticket_nombre"`
	FechaIni            time.Time
	FechaFin            time.Time
	Estado              int16 `gorm:"default:1"`
	UsuarioCreacion     *int64
	FechaCreacion       time.Time `gorm:"default:now()"`
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	Evento *Evento `gorm:"foreignKey:EventoID;references:evento_id"`
}

func (TipoDeTicket) TableName() string { return "tipo_de_ticket" }
