package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
)

type EventoFecha struct {
	ID                  uint64
	EventoID            uint64
	FechaID             uint64
	HoraInicio          time.Time
	Estado              util.Estado
	UsuarioCreacion     *uint64
	FechaCreacion       time.Time
	UsuarioModificacion *uint64
	FechaModificacion   *time.Time

	Tickets []Ticket
}
