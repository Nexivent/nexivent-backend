package domain

type Ticket struct {
	ID             int64          `db:"ticket_id"          json:"ticketId"`
	OrdenDeCompra  *OrdenDeCompra `db:"-"                  json:"ordenDeCompra,omitempty"` // FK opcional
	EventoFecha    EventoFecha    `db:"-"                  json:"eventoFecha"`             // FK -> evento_fecha (siempre)
	Tarifa         Tarifa         `db:"-"                  json:"tarifa"`                  // FK -> tarifa (siempre)
	CodigoQR       string         `db:"codigo_qr"          json:"codigoQR"`
	EstadoDeTicket int16          `db:"estado_de_ticket"   json:"estadoDeTicket"`
}
