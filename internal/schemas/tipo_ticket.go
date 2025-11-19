package schemas

type TipoTicketTicketRequest struct {
	EventoID int64  `json:"idEvento"`
	Nombre   string `json:"nombre"`
	FechaIni string `json:"fechaIni"` // "2006-01-02"
	FechaFin string `json:"fechaFin"` // "2006-01-02"
	Estado   int16  `json:"estado"`
}

type TipoTicketUpdateRequest struct {
	Nombre   *string `json:"nombre,omitempty"`
	FechaIni *string `json:"fechaIni,omitempty"`
	FechaFin *string `json:"fechaFin,omitempty"`
	Estado   *int16  `json:"estado,omitempty"`
}

type TipoTicketTicketResponse struct {
	ID       int64  `json:"idTipoTicket"`
	EventoID int64  `json:"idEvento"`
	Nombre   string `json:"nombre"`
	FechaIni string `json:"fechaIni"`
	FechaFin string `json:"fechaFin"`
	Estado   int16  `json:"estado"`
}
