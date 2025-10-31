package data

import (
	"time"
)

type Evento struct {
	ID                  uint64      `json:"id"`
	Nombre              string      `json:"nombre"`
	Organizador         Organizador `json:"organizador"`
	FechaCreacion       time.Time   `json:"-"`
	FechaRealizacion    time.Time   `json:"fecha_realizacion"`
	CantidadAsientos    uint32      `json:"cantidad_asientos"`
	AsientosDisponibles uint32      `json:"asientos_disponibles"`
	AsientosOcupados    uint32      `json:"asientos_ocupados"`
}


