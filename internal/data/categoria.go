package data

import (
	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/google/uuid"
)

type Categoria struct {
	ID          uuid.UUID
	Nombre      string
	Descripcion string
	Estado      util.Estado

	Eventos []Evento
}

func (Categoria) TableName() string { return "categoria" }
