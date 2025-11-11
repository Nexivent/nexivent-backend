package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
)

type Notificacion struct {
	ID                 uint64
	Mensaje            string
	Canal              string
	FechaEnvio         time.Time
	EstadoNotificacion util.EstadoNotificacion
}

func (Notificacion) TableName() string { return "notificacion" }
