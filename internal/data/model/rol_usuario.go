package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type RolUsuario struct {
	ID                  uint64      `gorm:"column:rol_usuario_id;primaryKey" json:"id"`
	RolID               uint64      `gorm:"column:rol_id" json:"rolId"`
	UsuarioID           uint64      `gorm:"column:usuario_id" json:"usuarioId"`
	UsuarioCreacion     *uint64     `gorm:"column:usuario_creacion" json:"-"`
	FechaCreacion       time.Time   `gorm:"column:fecha_creacion" json:"-"`
	UsuarioModificacion *uint64     `gorm:"column:usuario_modificacion" json:"-"`
	FechaModificacion   *time.Time  `gorm:"column:fecha_modificacion" json:"-"`
	Estado              util.Estado `gorm:"column:estado" json:"-"`
}

func (RolUsuario) TableName() string { return "rol_usuario" }

func ValidateRolUsuario(v *validator.Validator, rolUsuario *RolUsuario) {
	// Validar IDs
	v.Check(rolUsuario.RolID != 0, "rolId", "el ID del rol es obligatorio")
	v.Check(rolUsuario.UsuarioID != 0, "usuarioId", "el ID del usuario es obligatorio")
}
