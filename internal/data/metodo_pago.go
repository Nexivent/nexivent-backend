package data

import (
	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/google/uuid"
)

type MetodoDePago struct {
	ID     uuid.UUID
	Tipo   string
	Estado util.Estado

	Ordenes []OrdenDeCompra
}


