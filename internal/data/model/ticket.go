package model

type Ticket struct {
	ID              uint64
	OrdenDeCompraID *uint64
	EventoFechaID   uint64
	TarifaID        uint64
	CodigoQR        string
	EstadoDeTicket  int16
}

func (Ticket) TableName() string { return "ticket" }
