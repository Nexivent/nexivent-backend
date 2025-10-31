package domain

type Ticket struct {
	IDTicket              int            `db:"id_ticket" json:"idTicket"`
	OrdenDeCompraDeTicket OrdenDeCompra  `db:"-" json:"ordenDeCompra"`
	Evento                Evento         `db:"-" json:"evento"`
	CodigoQR              string         `db:"codigo_qr" json:"codigoQR"`
	Tipo                  TipoDeTicket   `db:"tipo"` // enum/string
	PrecioVenta           float64        `db:"precio_venta" json:"precioVenta"`
	Estado                EstadoDeTicket `db:"estado" json:"estado"`
}
