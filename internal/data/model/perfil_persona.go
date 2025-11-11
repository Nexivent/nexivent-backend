package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
)

type PerfilDePersona struct {
	ID                  uint64
	EventoID            uint64
	Nombre              string
	Estado              util.Estado
	UsuarioCreacion     uint64
	FechaCreacion       time.Time
	UsuarioModificacion uint64
	FechaModificacion   time.Time
}

func (PerfilDePersona) TableName() string { return "perfil_de_persona" }
