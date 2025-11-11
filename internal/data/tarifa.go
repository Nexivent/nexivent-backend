package data

import (
	"time"

	"github.com/google/uuid"
)

type Tarifa struct {
	ID                  uuid.UUID
	SectorID            uuid.UUID
	TipoDeTicketID      uuid.UUID
	PerfilDePersonaID   *uuid.UUID
	Precio              float64
	Estado              int16
	UsuarioCreacion     *uuid.UUID
	FechaCreacion       time.Time
	UsuarioModificacion *uuid.UUID
	FechaModificacion   *time.Time

	Sector        *Sector
	TipoDeTicket  *TipoDeTicket
	PerfilPersona *PerfilDePersona
}


