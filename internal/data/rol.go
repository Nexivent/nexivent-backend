package data

import (
	"time"

	"github.com/google/uuid"
)

type Rol struct {
	ID                  uuid.UUID
	Nombre              string
	UsuarioCreacion     *uuid.UUID
	FechaCreacion       time.Time
	UsuarioModificacion *uuid.UUID
	FechaModificacion   *time.Time

	Usuarios []RolUsuario
}

func (Rol) TableName() string { return "rol" }
