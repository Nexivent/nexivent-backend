package domain

import (
	"database/sql"
	"time"
)

type ComprobanteDePago struct {
	ID                int64           `db:"comprobante_de_pago_id" json:"id"`
	OrdenDeCompra     OrdenDeCompra   `db:"-" json:"ordenDeCompra"` // FK -> orden_de_compra
	TipoDeComprobante TipoComprobante `db:"tipo_de_comprobante" `
	Numero            string          `db:"numero" json:"numero"`
	FechaEmision      time.Time       `db:"fecha_emision" json:"fechaEmision"`
	RUC               sql.NullString  `db:"ruc" json:"ruc,omitempty"`
	DireccionFiscal   sql.NullString  `db:"direccion_fiscal" json:"direccionFiscal,omitempty"`
}
