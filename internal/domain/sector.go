package domain

import (
	"time"
)

type Sector struct {
	ID                  int64 `gorm:"column:sector_id;primaryKey;autoIncrement"`
	EventoID            int64
	SectorTipo          string
	TotalEntradas       int
	CantVendidas        int
	Estado              int16
	UsuarioCreacion     *int64
	FechaCreacion       time.Time
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	Evento *Evento `gorm:"foreignKey:EventoID"`
}

func (Sector) TableName() string { return "sector" }
