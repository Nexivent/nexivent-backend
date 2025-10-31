package domain

import "time"

type OrdenDeCompra struct {
	IDOrden               int        `db:"id_orden" json:"idOrden"`
	UsuarioComprador      Usuario    `db:"-" json:"usuarioComprador"`
	FechaDeCompra         time.Time  `db:"fecha_compra" json:"fechaDeCompra"`
	Total                 float64    `db:"total" json:"total"`
	MontoFeeServicio      float64    `db:"monto_fee_servicio" json:"montoFeeServicio"`
	Estado                string     `db:"estado" json:"estado"`
	MetodoDePago          MetodoPago `db:"-" json:"metodoPago,omitempty"`
	IDTransaccionExterna  string     `db:"id_transaccion_externa" json:"idTransaccionExterna"`
	FechaLiberacionFondos time.Time  `db:"fecha_liberacion_fondos" json:"fechaLiberacionFondos"`
}
