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
	Estado              int16 `gorm:"default:1"`
	UsuarioCreacion     *int64
	FechaCreacion       time.Time `gorm:"default:now()"`
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	Sector        *Sector          `gorm:"foreignKey:SectorID;references:sector_id"`
	TipoDeTicket  *TipoDeTicket    `gorm:"foreignKey:TipoDeTicketID;references:tipo_de_ticket_id"`
	PerfilPersona *PerfilDePersona `gorm:"foreignKey:PerfilDePersonaID;references:perfil_de_persona_id"`
}

func (Tarifa) TableName() string { return "tarifa" }
