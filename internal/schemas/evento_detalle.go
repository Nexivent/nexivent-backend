package schemas

// DTOs para respuesta JSON
type FechaEventoDTO struct {
	IDFechaEvento int64  `json:"idFechaEvento"`
	Fecha         string `json:"fecha"`
	HoraInicio    string `json:"horaInicio"`
	HoraFin       string `json:"horaFin"`
}

type TarifaDTO struct {
	IDTarifa        int64   `json:"idTarifa"`
	Precio          float64 `json:"precio"`
	TipoSector      string  `json:"tipoSector"`
	StockDisponible int     `json:"stockDisponible"`
	TipoTicket      string  `json:"tipoTicket"`
	FechaIni        string  `json:"fechaIni"`
	FechaFin        string  `json:"fechaFin"`
	Perfil          string  `json:"perfil"`
}

type EventoDetalleDTO struct {
	IDEvento    int64            `json:"idEvento"`
	Titulo      string           `json:"titulo"`
	Descripcion string           `json:"descripcion"`
	ImagenPortada string         `json:"imagenPortada"`
	Lugar       string           `json:"lugar"`
	Fechas      []FechaEventoDTO `json:"fechas"`
	Tarifas     []TarifaDTO      `json:"tarifas"`
}