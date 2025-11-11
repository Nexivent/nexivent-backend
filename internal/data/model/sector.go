package model

import (
	"time"
)

type Sector struct {
	ID                  uint64 
	EventoID            uint64 
	SectorTipo          string 
	TotalEntradas       int
	CantVendidas        int  
	Estado              int16 
	UsuarioCreacion     *uint64
	FechaCreacion       time.Time 
	UsuarioModificacion *uint64
	FechaModificacion   *time.Time
}

func (Sector) TableName() string { return "sector" }
