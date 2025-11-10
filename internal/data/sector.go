package data

import (
	"time"

	"github.com/google/uuid"
)

type Sector struct {
	ID                  uuid.UUID `gorm:"column:sector_id;primaryKey"`
	EventoID            uuid.UUID `gorm:"uniqueIndex:uq_sector_tipo"`
	SectorTipo          string    `gorm:"uniqueIndex:uq_sector_tipo"`
	TotalEntradas       int
	CantVendidas        int   `gorm:"default:0"`
	Estado              int16 `gorm:"default:1"`
	UsuarioCreacion     *uuid.UUID
	FechaCreacion       time.Time `gorm:"default:now()"`
	UsuarioModificacion *uuid.UUID
	FechaModificacion   *time.Time

	Evento *Evento `gorm:"foreignKey:EventoID;references:evento_id"`
}

func (Sector) TableName() string { return "sector" }
