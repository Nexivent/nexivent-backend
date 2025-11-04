package domain

import (
	"time"
)

type ComprobanteDePago struct {
	ID                int64 `gorm:"column:comprobante_de_pago_id;primaryKey;autoIncrement"`
	OrdenDeCompraID   int64
	TipoDeComprobante int16
	Numero            string
	FechaEmision      time.Time
	RUC               *string
	DireccionFiscal   *string

	OrdenDeCompra *OrdenDeCompra `gorm:"foreignKey:OrdenDeCompraID"`
}

func (ComprobanteDePago) TableName() string { return "comprobante_de_pago" }
