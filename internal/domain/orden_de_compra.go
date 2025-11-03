package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrdenDeCompra struct {
	ID               int64             `db:"orden_de_compra_id" json:"id"`
	Usuario          Usuario           `db:"-" json:"usuario"`      // FK -> usuario
	MetodoDePago     MetodoDePago      `db:"-" json:"metodoDePago"` // FK -> metodo_de_pago
	Fecha            time.Time         `db:"fecha" json:"fecha"`
	Total            decimal.Decimal   `db:"total" json:"total"`
	MontoFeeServicio decimal.Decimal   `db:"monto_fee_servicio" json:"montoFeeServicio"`
	EstadoDeOrden    EstadoOrdenCompra `db:"estado_de_orden"`
}
