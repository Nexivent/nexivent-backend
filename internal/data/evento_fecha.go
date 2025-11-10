package data

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/google/uuid"
)

type EventoFecha struct {
	ID                  uuid.UUID
	EventoID            uuid.UUID
	FechaID             uuid.UUID
	HoraInicio          time.Time
	Estado              util.Estado
	UsuarioCreacion     *uuid.UUID
	FechaCreacion       time.Time
	UsuarioModificacion *uuid.UUID
	FechaModificacion   *time.Time

	Evento *Evento
	Fecha  *Fecha

	Tickets []Ticket
}

func (EventoFecha) TableName() string { return "evento_fecha" }
