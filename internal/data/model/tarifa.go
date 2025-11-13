package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type Tarifa struct {
	ID                  uint64      `gorm:"column:tarifa_id;primaryKey" json:"id"`
	SectorID            uint64      `gorm:"column:sector_id" json:"sectorId"`
	TipoDeTicketID      uint64      `gorm:"column:tipo_de_ticket_id" json:"tipoDeTicketId"`
	PerfilDePersonaID   *uint64     `gorm:"column:perfil_de_persona_id" json:"perfilDePersonaId,omitempty"`
	Precio              float64     `gorm:"column:precio" json:"precio"`
	Estado              util.Estado `gorm:"column:estado" json:"-"`
	UsuarioCreacion     *uint64     `gorm:"column:usuario_creacion" json:"-"`
	FechaCreacion       time.Time   `gorm:"column:fecha_creacion" json:"-"`
	UsuarioModificacion *uint64     `gorm:"column:usuario_modificacion" json:"-"`
	FechaModificacion   *time.Time  `gorm:"column:fecha_modificacion" json:"-"`

	Tickets []Ticket `json:"tickets,omitempty"`
}

func (Tarifa) TableName() string { return "tarifa" }

func ValidateTarifa(v *validator.Validator, tarifa *Tarifa) {
	// Validar IDs
	v.Check(tarifa.SectorID != 0, "sectorId", "el ID del sector es obligatorio")
	v.Check(tarifa.TipoDeTicketID != 0, "tipoDeTicketId", "el ID del tipo de ticket es obligatorio")

	// Validar Precio
	v.Check(tarifa.Precio >= 0, "precio", "el precio no puede ser negativo")
}
