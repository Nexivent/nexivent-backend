package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
)

type ComprobanteDePago struct {
	ID                uint64
	OrdenDeCompraID   uint64
	TipoDeComprobante util.TipoComprobante
	Numero            string
	FechaEmision      time.Time
	RUC               *string
	DireccionFiscal   *string
}

func (ComprobanteDePago) TableName() string { return "comprobante_de_pago" }