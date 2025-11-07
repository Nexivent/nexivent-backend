package model

import (
	"time"
)

type RolUsuario struct {
	ID                  int64 `gorm:"column:rol_usuario_id;primaryKey;autoIncrement"`
	RolID               int64 `gorm:"uniqueIndex:ux_usuario_rol"`
	UsuarioID           int64 `gorm:"uniqueIndex:ux_usuario_rol"`
	UsuarioCreacion     int64
	FechaCreacion       time.Time
	UsuarioModificacion *int64
	FechaModificacion   *time.Time
	Estado              int16

	Rol     *Rol     `gorm:"foreignKey:RolID"`
	Usuario *Usuario `gorm:"foreignKey:UsuarioID"`
}

func (RolUsuario) TableName() string { return "rol_usuario" }
