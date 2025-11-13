package model

import (
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type UsuarioCupon struct {
	CuponID   uint64 `gorm:"column:cupon_id;primaryKey" json:"cuponId"`
	UsuarioID uint64 `gorm:"column:usuario_id;primaryKey" json:"usuarioId"`
	CantUsada int64  `gorm:"column:cant_usada" json:"cantUsada"`
}

func (UsuarioCupon) TableName() string { return "usuario_cupon" }

func ValidateUsuarioCupon(v *validator.Validator, usuarioCupon *UsuarioCupon) {
	// Validar IDs
	v.Check(usuarioCupon.CuponID != 0, "cuponId", "el ID del cupón es obligatorio")
	v.Check(usuarioCupon.UsuarioID != 0, "usuarioId", "el ID del usuario es obligatorio")

	// Validar CantUsada
	v.Check(usuarioCupon.CantUsada >= 0, "cantUsada", "la cantidad usada no puede ser negativa")
}
