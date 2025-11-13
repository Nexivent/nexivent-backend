package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type TipoDeTicket struct {
	ID                  uint64      `gorm:"column:tipo_de_ticket_id;primaryKey" json:"id"`
	EventoID            uint64      `gorm:"column:evento_id" json:"eventoId"`
	Nombre              string      `gorm:"column:nombre" json:"nombre"`
	FechaIni            time.Time   `gorm:"column:fecha_ini" json:"fechaIni"`
	FechaFin            time.Time   `gorm:"column:fecha_fin" json:"fechaFin"`
	Estado              util.Estado `gorm:"column:estado" json:"-"`
	UsuarioCreacion     *uint64     `gorm:"column:usuario_creacion" json:"-"`
	FechaCreacion       time.Time   `gorm:"column:fecha_creacion" json:"-"`
	UsuarioModificacion *uint64     `gorm:"column:usuario_modificacion" json:"-"`
	FechaModificacion   *time.Time  `gorm:"column:fecha_modificacion" json:"-"`

	Tarifas []Tarifa `json:"tarifas,omitempty"`
}

func (TipoDeTicket) TableName() string { return "tipo_de_ticket" }

func ValidateTipoDeTicket(v *validator.Validator, tipoDeTicket *TipoDeTicket) {
	// Validar EventoID
	v.Check(tipoDeTicket.EventoID != 0, "eventoId", "el ID del evento es obligatorio")

	// Validar Nombre
	v.Check(tipoDeTicket.Nombre != "", "nombre", "el nombre es obligatorio")
	v.Check(len(tipoDeTicket.Nombre) <= 100, "nombre", "el nombre no debe exceder 100 caracteres")

	// Validar rango de fechas
	v.Check(!tipoDeTicket.FechaIni.IsZero(), "fechaIni", "la fecha de inicio es obligatoria")
	v.Check(!tipoDeTicket.FechaFin.IsZero(), "fechaFin", "la fecha de fin es obligatoria")
	v.Check(!tipoDeTicket.FechaFin.Before(tipoDeTicket.FechaIni), "fechaFin", "la fecha de fin debe ser mayor o igual a la fecha de inicio")
}
