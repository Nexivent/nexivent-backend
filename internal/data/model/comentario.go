package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
)

type Comentario struct {
	ID            uint64
	UsuarioID     uint64
	EventoID      uint64
	Descripcion   string
	FechaCreacion time.Time
	Estado        util.Estado
}

func (Comentario) TableName() string { return "comentario" }
