package data

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/google/uuid"
)

type TipoDeTicket struct {
	ID                  uuid.UUID
	EventoID            uuid.UUID
	Nombre              string
	FechaIni            time.Time
	FechaFin            time.Time
	Estado              util.Estado
	UsuarioCreacion     *uuid.UUID
	FechaCreacion       time.Time
	UsuarioModificacion *uuid.UUID
	FechaModificacion   *time.Time

	Evento *Evento
}


