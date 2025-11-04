package model

import (
	"time"
)

type Tarifa struct {
	ID                  int64 `gorm:"column:tarifa_id;primaryKey;autoIncrement"`
	SectorID            int64
	TipoDeTicketID      int64
	PerfilDePersonaID   *int64
	Precio              float64
	Estado              int16
	UsuarioCreacion     *int64
	FechaCreacion       time.Time
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	Sector        *Sector          `gorm:"foreignKey:SectorID"`
	TipoDeTicket  *TipoDeTicket    `gorm:"foreignKey:TipoDeTicketID"`
	PerfilPersona *PerfilDePersona `gorm:"foreignKey:PerfilDePersonaID"`
}

func (Tarifa) TableName() string { return "tarifa" }
