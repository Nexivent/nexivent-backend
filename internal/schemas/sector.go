package schemas

type SectorTicketRequest struct {
	EventoID      int64  `json:"idEvento"`
	SectorTipo    string `json:"sector"`
	TotalEntradas int    `json:"totalEntradas"`
	Estado        int16  `json:"estado"` // 1 por defecto
}

type SectorUpdateRequest struct {
	SectorTipo    *string `json:"sector,omitempty"`
	TotalEntradas *int    `json:"totalEntradas,omitempty"`
	CantVendidas  *int    `json:"cantVendidas,omitempty"`
	Estado        *int16  `json:"estado,omitempty"`
}

type SectorTicketResponse struct {
	ID            int64  `json:"idSector"`
	EventoID      int64  `json:"idEvento"`
	SectorTipo    string `json:"sector"`
	TotalEntradas int    `json:"totalEntradas"`
	CantVendidas  int    `json:"cantVendidas"`
	Estado        int16  `json:"estado"`
}
