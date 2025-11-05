package model

import (
	"time"
)

type Rol struct {
	ID                  int64 `gorm:"column:rol_id;primaryKey;autoIncrement"`
	Nombre              string
	UsuarioCreacion     int64
	FechaCreacion       time.Time
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	Usuarios []RolUsuario
}

func (Rol) TableName() string { return "rol" }
