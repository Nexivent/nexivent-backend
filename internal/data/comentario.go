package data

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/google/uuid"
)

type Comentario struct {
	ID            uuid.UUID
	UsuarioID     uuid.UUID
	EventoID      uuid.UUID
	Descripcion   string
	FechaCreacion time.Time
	Estado        util.Estado

	Usuario *Usuario
	Evento  *Evento
}

