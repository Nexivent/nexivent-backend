package model

import (
	"time"
)

type Rol struct {
	ID                  uint64
	Nombre              string
	UsuarioCreacion     *uint64
	FechaCreacion       time.Time
	UsuarioModificacion *uint64
	FechaModificacion   *time.Time

	Usuarios []RolUsuario
}

func (Rol) TableName() string { return "rol" }
