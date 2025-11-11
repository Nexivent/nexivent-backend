package model

import (
	"time"
)

type Tarifa struct {
	ID                  uint64
	SectorID            uint64
	TipoDeTicketID      uint64
	PerfilDePersonaID   *uint64
	Precio              float64
	Estado              int16
	UsuarioCreacion     *uint64
	FechaCreacion       time.Time
	UsuarioModificacion *uint64
	FechaModificacion   *time.Time
}

func (Tarifa) TableName() string { return "tarifa" }
