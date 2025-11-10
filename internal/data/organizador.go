package data

import (
	"time"

	"github.com/google/uuid"
)

type Organizador struct {
	RUC           uuid.UUID `json:"ruc"`
	RazonSocial   string    `json:"razon_social"`
	FechaCreacion time.Time `json:"-"`
	Eventos       []Evento  `json:"eventos"`
}
