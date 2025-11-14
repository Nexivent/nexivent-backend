package model

import (
	"github.com/Nexivent/nexivent-backend/internal/data/model/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type Ticket struct {
	ID              uint64              `gorm:"column:ticket_id;primaryKey;autoIncrement" json:"id"`
	OrdenDeCompraID *uint64             `gorm:"column:orden_de_compra_id" json:"ordenDeCompraId,omitempty"`
	EventoFechaID   uint64              `gorm:"column:evento_fecha_id" json:"eventoFechaId"`
	TarifaID        uint64              `gorm:"column:tarifa_id" json:"tarifaId"`
	CodigoQR        string              `gorm:"column:codigo_qr;uniqueIndex" json:"codigoQr"`
	EstadoDeTicket  util.EstadoDeTicket `gorm:"column:estado_de_ticket;default:0" json:"estadoDeTicket"`

	// OrdenDeCompra *OrdenDeCompra `gorm:"foreignKey:OrdenDeCompraID;references:orden_de_compra_id"`
	// EventoFecha   *EventoFecha   `gorm:"foreignKey:EventoFechaID;references:evento_fecha_id"`
	// Tarifa        *Tarifa        `gorm:"foreignKey:TarifaID;references:tarifa_id"`
}

func (Ticket) TableName() string { return "ticket" }

func ValidateTicket(v *validator.Validator, ticket *Ticket) {
	// Validar IDs
	v.Check(ticket.EventoFechaID != 0, "eventoFechaId", "el ID de la fecha del evento es obligatorio")
	v.Check(ticket.TarifaID != 0, "tarifaId", "el ID de la tarifa es obligatorio")

	// Validar CodigoQR
	v.Check(ticket.CodigoQR != "", "codigoQr", "el código QR es obligatorio")
	v.Check(len(ticket.CodigoQR) <= 255, "codigoQr", "el código QR no debe exceder 255 caracteres")
}
