package domain

type Ticket struct {
	ID             int64          `db:"ticket_id" json:"id"`
	OrdenDeCompra  *OrdenDeCompra `db:"-" json:"ordenDeCompra,omitempty"` // FK opcional
	EventoFecha    EventoFecha    `db:"-" json:"eventoFecha"`             // FK -> evento_fecha
	Tarifa         Tarifa         `db:"-" json:"tarifa"`                  // FK -> tarifa
	CodigoQR       string         `db:"codigo_qr" json:"codigoQR"`
	EstadoDeTicket int16          `db:"estado_de_ticket"`
}
