package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrdenDeCompra struct {
	ID               int64           `db:"orden_de_compra_id" json:"ordenDeCompraId"`
	Usuario          Usuario         `db:"-"                   json:"usuario"`      // FK -> usuario (siempre)
	MetodoDePago     MetodoDePago    `db:"-"                   json:"metodoDePago"` // FK -> metodo_de_pago (siempre)
	Fecha            time.Time       `db:"fecha"               json:"fecha"`        // DATE
	Total            decimal.Decimal `db:"total"               json:"total"`
	MontoFeeServicio decimal.Decimal `db:"monto_fee_servicio"  json:"montoFeeServicio"`
	EstadoDeOrden    int16           `db:"estado_de_orden"     json:"estadoDeOrden"`
}
