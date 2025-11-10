package data

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/google/uuid"
)

type RolUsuario struct {
	ID                  uuid.UUID
	RolID               uuid.UUID
	UsuarioID           uuid.UUID
	UsuarioCreacion     *uuid.UUID
	FechaCreacion       time.Time
	UsuarioModificacion *uuid.UUID
	FechaModificacion   *time.Time
	Estado              util.Estado

	Rol     *Rol
	Usuario *Usuario
}

func (RolUsuario) TableName() string { return "rol_usuario" }
