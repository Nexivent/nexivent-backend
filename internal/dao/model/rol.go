package model

import (
	"time"
)

type Rol struct {
	ID                  int64  `gorm:"column:rol_id;primaryKey;autoIncrement"`
	Nombre              string `gorm:"uniqueIndex"`
	UsuarioCreacion     *int64
	FechaCreacion       time.Time `gorm:"default:now()"`
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	Usuarios []RolUsuario
}

func (Rol) TableName() string { return "rol" }
