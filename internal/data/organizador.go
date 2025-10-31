package data

import (
	"time"
)

type Organizador struct {
	RUC           uint64    `json:"ruc"`
	RazonSocial   string    `json:"razon_social"`
	FechaCreacion time.Time `json:"-"`
	Eventos       []Evento  `json:"eventos"`
}
