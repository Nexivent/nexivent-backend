package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type Rol struct {
	ID                  uint64     `gorm:"column:rol_id;primaryKey" json:"id"`
	Nombre              string     `gorm:"column:nombre;uniqueIndex" json:"nombre"`
	UsuarioCreacion     *uint64    `gorm:"column:usuario_creacion" json:"-"`
	FechaCreacion       time.Time  `gorm:"column:fecha_creacion;default:now()" json:"-"`
	UsuarioModificacion *uint64    `gorm:"column:usuario_modificacion" json:"-"`
	FechaModificacion   *time.Time `gorm:"column:fecha_modificacion" json:"-"`

	Usuarios []RolUsuario `json:"usuarios"`
}

func (Rol) TableName() string { return "rol" }

func ValidateRol(v *validator.Validator, rol *Rol) {
	// Validar Nombre
	v.Check(rol.Nombre != "", "nombre", "el nombre es obligatorio")
	v.Check(len(rol.Nombre) <= 100, "nombre", "el nombre no debe exceder 100 caracteres")
}
