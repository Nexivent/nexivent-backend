package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
)

type Cupon struct {
	ID                  uint64
	Descripcion         string
	Tipo                string
	Valor               float64
	Estado              util.Estado
	Codigo              string
	UsoPorUsuario       int64
	UsoRealizados       int64
	UsuarioCreacion     *uint64
	FechaCreacion       time.Time
	UsuarioModificacion *uint64
	FechaModificacion   *time.Time
	EventoID uint64

	// Mantienes tu relación con usuarios
	Usuarios []UsuarioCupon
}
