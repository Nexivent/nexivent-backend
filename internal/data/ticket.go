package data

import "github.com/google/uuid"

type Ticket struct {
	ID              uuid.UUID
	OrdenDeCompraID *uuid.UUID
	EventoFechaID   uuid.UUID
	TarifaID        uuid.UUID
	CodigoQR        string
	EstadoDeTicket  int16

	OrdenDeCompra *OrdenDeCompra
	EventoFecha   *EventoFecha
	Tarifa        *Tarifa
}


