package model

import (
	"time"
)

type EventoFecha struct {
	ID                      int64     `gorm:"column:evento_fecha_id;primaryKey;autoIncrement"`
	EventoID                int64     `gorm:"uniqueIndex:uq_evento_fecha"`
	FechaID                 int64     `gorm:"uniqueIndex:uq_evento_fecha"`
	HoraInicio              time.Time `gorm:"type:time;uniqueIndex:uq_evento_fecha"`
	GananciaNetaOrganizador float64   `gorm:"column:ganancia_neta_organizador;default:0" json:"ganancia_neta_organizador"`
	Estado                  int16     `gorm:"default:1"`
	UsuarioCreacion         *int64
	FechaCreacion           time.Time `gorm:"default:now()"`
	UsuarioModificacion     *int64
	FechaModificacion       *time.Time

	Evento *Evento `gorm:"foreignKey:EventoID;references:evento_id"`
	Fecha  *Fecha  `gorm:"foreignKey:FechaID;references:fecha_id"`

	Tickets []Ticket
}

func (EventoFecha) TableName() string { return "evento_fecha" }
