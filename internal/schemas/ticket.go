package schemas

// Request para emitir tickets a partir de una orden confirmada
type TicketIssueRequest struct {
	OrderID int64 `json:"orderId"`
}

// Ticket emitido hacia el front
type TicketEmitido struct {
	IdTicket     int64  `json:"idTicket"`
	CodigoQR     string `json:"codigoQR"`
	Estado       string `json:"estado"`       // "DISPONIBLE"
	TituloEvento string `json:"tituloEvento"` // título del evento
	FechaEvento  string `json:"fechaEvento"`  // YYYY-MM-DD
	HoraInicio   string `json:"horaInicio"`   // HH:MM
	Sector       string `json:"sector"`       // descripción del sector
}

// Response de emisión de tickets
type TicketIssueResponse struct {
	Tickets []TicketEmitido `json:"tickets"`
}

// Request para cancelar tickets
type TicketCancelRequest struct {
	IdTickets []int64 `json:"idTickets"`
}

// Ticket cancelado
type TicketCancelado struct {
	IdTicket int64  `json:"idTicket"`
	Estado   string `json:"estado"` // "CANCELADO"
}

// Response de cancelación
type TicketCancelResponse struct {
	Cancelados    []TicketCancelado `json:"cancelados"`
	NoEncontrados []int64           `json:"noEncontrados,omitempty"`
	NoCancelables []int64           `json:"noCancelables,omitempty"`
	Mensaje       string            `json:"mensaje"`
}
