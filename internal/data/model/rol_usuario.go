package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
)

type RolUsuario struct {
	ID                  uint64
	RolID               uint64
	UsuarioID           uint64
	UsuarioCreacion     *uint64
	FechaCreacion       time.Time
	UsuarioModificacion *uint64
	FechaModificacion   *time.Time
	Estado              util.Estado
}

func (RolUsuario) TableName() string { return "rol_usuario" }
