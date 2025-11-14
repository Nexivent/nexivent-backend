package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/model/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type Sector struct {
	ID                  uint64      `gorm:"column:sector_id;primaryKey;autoIncrement" json:"id"`
	EventoID            uint64      `gorm:"column:evento_id;uniqueIndex:uq_sector_tipo" json:"eventoId"`
	SectorTipo          string      `gorm:"column:sector_tipo;uniqueIndex:uq_sector_tipo" json:"sectorTipo"`
	TotalEntradas       int         `gorm:"column:total_entradas" json:"totalEntradas"`
	CantVendidas        int         `gorm:"column:cant_vendidas;default:0" json:"cantVendidas"`
	Estado              util.Estado `gorm:"column:estado;default:1" json:"-"`
	UsuarioCreacion     *uint64     `gorm:"column:usuario_creacion" json:"-"`
	FechaCreacion       time.Time   `gorm:"column:fecha_creacion;default:now()" json:"-"`
	UsuarioModificacion *uint64     `gorm:"column:usuario_modificacion" json:"-"`
	FechaModificacion   *time.Time  `gorm:"column:fecha_modificacion" json:"-"`

	// Evento *Evento `gorm:"foreignKey:EventoID;references:evento_id"`
}

func (Sector) TableName() string { return "sector" }

func ValidateSector(v *validator.Validator, sector *Sector) {
	// Validar EventoID
	v.Check(sector.EventoID != 0, "eventoId", "el ID del evento es obligatorio")

	// Validar SectorTipo
	v.Check(sector.SectorTipo != "", "sectorTipo", "el tipo de sector es obligatorio")
	v.Check(len(sector.SectorTipo) <= 50, "sectorTipo", "el tipo de sector no debe exceder 50 caracteres")

	// Validar TotalEntradas
	v.Check(sector.TotalEntradas > 0, "totalEntradas", "el total de entradas debe ser mayor a 0")

	// Validar CantVendidas
	v.Check(sector.CantVendidas >= 0, "cantVendidas", "la cantidad vendida no puede ser negativa")
	v.Check(sector.CantVendidas <= sector.TotalEntradas, "cantVendidas", "la cantidad vendida no puede superar el total de entradas")
}
