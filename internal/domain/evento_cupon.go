package domain

import (
	"time"
)

type EventoCupon struct {
	EventoID            int64 `gorm:"primaryKey"`
	CuponID             int64 `gorm:"primaryKey"`
	CantCupones         int64
	FechaIni            time.Time
	FechaFin            time.Time
	UsuarioCreacion     *int64
	FechaCreacion       time.Time
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	Evento *Evento `gorm:"foreignKey:EventoID"`
	Cupon  *Cupon  `gorm:"foreignKey:CuponID"`
}

func (EventoCupon) TableName() string { return "evento_cupon" }
