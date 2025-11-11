package model

import (
	"github.com/Nexivent/nexivent-backend/internal/data/util"
)

type MetodoDePago struct {
	ID     uint64
	Tipo   string
	Estado util.Estado

	Ordenes []OrdenDeCompra
}

func (MetodoDePago) TableName() string { return "metodo_de_pago" }
