package schemas

type TarifaRequest struct {
	SectorID          int64   `json:"idSector"`
	TipoDeTicketID    int64   `json:"idTipoTicket"`
	PerfilDePersonaID *int64  `json:"idPerfilPersona,omitempty"`
	Precio            float64 `json:"precio"`
	Estado            int16   `json:"estado"` // 1
}

type TarifaUpdateRequest struct {
	SectorID          *int64   `json:"idSector,omitempty"`
	TipoDeTicketID    *int64   `json:"idTipoTicket,omitempty"`
	PerfilDePersonaID *int64   `json:"idPerfilPersona,omitempty"`
	Precio            *float64 `json:"precio,omitempty"`
	Estado            *int16   `json:"estado,omitempty"`
}

type TarifaResponse struct {
	ID                int64   `json:"idTarifa"`
	SectorID          int64   `json:"idSector"`
	TipoDeTicketID    int64   `json:"idTipoTicket"`
	PerfilDePersonaID *int64  `json:"idPerfilPersona,omitempty"`
	Precio            float64 `json:"precio"`
	Estado            int16   `json:"estado"`
}

type TarifaResponseOtros struct {
	ID     int64   `json:"idTarifa"`
	Precio float64 `json:"precio"`
	Estado int16   `json:"estado"`
}
