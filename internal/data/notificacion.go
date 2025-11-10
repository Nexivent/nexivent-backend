package data

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/google/uuid"
)

type Notificacion struct {
	ID                 uuid.UUID
	Mensaje            string
	Canal              string
	FechaEnvio         time.Time
	EstadoNotificacion util.EstadoNotificacion
}

func (Notificacion) TableName() string { return "notificacion" }
