package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
)

type TipoDeTicket struct {
	ID                  uint64
	EventoID            uint64
	Nombre              string
	FechaIni            time.Time
	FechaFin            time.Time
	Estado              util.Estado
	UsuarioCreacion     *uint64
	FechaCreacion       time.Time
	UsuarioModificacion *uint64
	FechaModificacion   *time.Time
}

func (TipoDeTicket) TableName() string { return "tipo_de_ticket" }
