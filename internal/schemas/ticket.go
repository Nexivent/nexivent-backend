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

// Request para emitir tickets CON INFO del frontend
type EmitirTicketsRequest struct {
	OrderID       int64               `json:"orderId"`
	UserID        int64               `json:"userId"`
	IdEvento      int64               `json:"idEvento"`
	IdFechaEvento int64               `json:"idFechaEvento"`
	Tickets       []TicketEmisionInfo `json:"tickets"`
}

type TicketEmisionInfo struct {
	IdTarifa     int64   `json:"idTarifa"`
	IdSector     int64   `json:"idSector"`
	IdPerfil     int64   `json:"idPerfil"`
	IdTipoTicket int64   `json:"idTipoTicket"`
	Cantidad     int     `json:"cantidad"`
	Precio       float64 `json:"precio"`
	NombreZona   string  `json:"nombreZona"`
}

// Response con los tickets generados
type EmitirTicketsResponse struct {
	Tickets []TicketGenerado `json:"tickets"`
	OrderID int64            `json:"orderId"`
}

type TicketGenerado struct {
	IdTicket string `json:"idTicket"`
	CodigoQR string `json:"codigoQR"`
	Estado   string `json:"estado"`
	Zona     string `json:"zona"`
}


type TicketDetalle struct {
	IDTicket    int64     `json:"idTicket"`
	TipoSector  string    `json:"tipoSector"`
	Evento      EventoMini    `json:"evento"`
	FechaInicio string `json:"fechaInicio"`
}
type EventoMini struct {
	IDEvento      int64  `json:"idEvento"`
	Titulo        string `json:"titulo"`
	Lugar         string `json:"lugar"`
	ImagenPortada string `json:"imagenPortada"`
}