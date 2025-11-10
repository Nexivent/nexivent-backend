package model

import (
	"time"
)

type RolUsuario struct {
	ID                  int64 `gorm:"column:rol_usuario_id;primaryKey;autoIncrement"`
	RolID               int64 `gorm:"uniqueIndex:uq_usuario_rol"`
	UsuarioID           int64 `gorm:"uniqueIndex:uq_usuario_rol"`
	UsuarioCreacion     *int64
	FechaCreacion       time.Time `gorm:"default:now()"`
	UsuarioModificacion *int64
	FechaModificacion   *time.Time
	Estado              int16 `gorm:"default:1"`

	Rol     *Rol     `gorm:"foreignKey:RolID;references:rol_id"`
	Usuario *Usuario `gorm:"foreignKey:UsuarioID;references:usuario_id"`
}

func (RolUsuario) TableName() string { return "rol_usuario" }
