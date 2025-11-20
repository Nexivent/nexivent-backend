package model

import (
	"time"
)

type Sector struct {
	ID                  int64  `gorm:"column:sector_id;primaryKey;autoIncrement"`
	EventoID            int64  `gorm:"uniqueIndex:uq_sector_tipo"`
	SectorTipo          string `gorm:"uniqueIndex:uq_sector_tipo"`
	TotalEntradas       int
	CantVendidas        int   `gorm:"default:0"`
	Estado              int16 `gorm:"default:1"`
	UsuarioCreacion     *int64
	FechaCreacion       time.Time `gorm:"default:now()"`
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	Tarifa []Tarifa `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Evento *Evento  `gorm:"foreignKey:EventoID;references:evento_id"`
}

func (Sector) TableName() string { return "sector" }
