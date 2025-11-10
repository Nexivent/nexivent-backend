package data

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/google/uuid"
)

type Cupon struct {
	ID                  uuid.UUID
	Descripcion         string
	Tipo                string
	Valor               float64
	Estado              util.Estado
	Codigo              string
	UsoPorUsuario       int64
	UsoRealizados       int64
	UsuarioCreacion     *uuid.UUID
	FechaCreacion       time.Time
	UsuarioModificacion *uuid.UUID
	FechaModificacion   *time.Time

	// FK al evento (muchos cupones pertenecen a un evento)
	EventoID uuid.UUID
	Evento   *Evento

	// Mantienes tu relación con usuarios
	Usuarios []UsuarioCupon
}

func (Cupon) TableName() string { return "cupon" }
