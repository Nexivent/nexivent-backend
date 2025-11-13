package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type PerfilDePersona struct {
	ID                  uint64      `gorm:"column:perfil_de_persona_id;primaryKey" json:"id"`
	EventoID            uint64      `gorm:"column:evento_id" json:"eventoId"`
	Nombre              string      `gorm:"column:nombre" json:"nombre"`
	Estado              util.Estado `gorm:"column:estado" json:"-"`
	UsuarioCreacion     *uint64     `gorm:"column:usuario_creacion" json:"-"`
	FechaCreacion       time.Time   `gorm:"column:fecha_creacion" json:"-"`
	UsuarioModificacion *uint64     `gorm:"column:usuario_modificacion" json:"-"`
	FechaModificacion   *time.Time  `gorm:"column:fecha_modificacion" json:"-"`

	Tarifas []Tarifa `json:"tarifas,omitempty"`
}

func (PerfilDePersona) TableName() string { return "perfil_de_persona" }

func ValidatePerfilDePersona(v *validator.Validator, perfil *PerfilDePersona) {
	// Validar EventoID
	v.Check(perfil.EventoID != 0, "eventoId", "el ID del evento es obligatorio")

	// Validar Nombre
	v.Check(perfil.Nombre != "", "nombre", "el nombre es obligatorio")
	v.Check(len(perfil.Nombre) <= 100, "nombre", "el nombre no debe exceder 100 caracteres")
}
