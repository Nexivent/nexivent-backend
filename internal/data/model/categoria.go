package model

import (
	"github.com/Nexivent/nexivent-backend/internal/data/util"
)

type Categoria struct {
	ID          uint64  
	Nombre      string 
	Descripcion string 
	Estado      util.Estado

	Eventos []Evento
}

func (Categoria) TableName() string { return "categoria" }