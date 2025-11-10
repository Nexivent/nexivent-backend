package data

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/google/uuid"
)

type ComprobanteDePago struct {
	ID                uuid.UUID
	OrdenDeCompraID   uuid.UUID
	TipoDeComprobante util.TipoComprobante
	Numero            string
	FechaEmision      time.Time
	RUC               *string
	DireccionFiscal   *string

	OrdenDeCompra *OrdenDeCompra
}

func (ComprobanteDePago) TableName() string { return "comprobante_de_pago" }
