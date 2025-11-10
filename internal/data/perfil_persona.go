package data

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/google/uuid"
)

type PerfilDePersona struct {
	ID                  uuid.UUID
	EventoID            uuid.UUID
	Nombre              string
	Estado              util.Estado
	UsuarioCreacion     uuid.UUID
	FechaCreacion       time.Time
	UsuarioModificacion uuid.UUID
	FechaModificacion   time.Time

	Evento Evento
}

func (PerfilDePersona) TableName() string { return "perfil_de_persona" }
