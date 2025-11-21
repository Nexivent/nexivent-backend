package schemas

// AdminReportRequest: Filtros de entrada
type AdminReportRequest struct {
	FechaInicio   *string `json:"fechaInicio"` // ISO
	FechaFin      *string `json:"fechaFin"`    // ISO
	IdCategoria   *int64  `json:"idCategoria"`
	IdOrganizador *int64  `json:"idOrganizador"`
	Estado        string  `json:"estado"` // PUBLICADO | CANCELADO | BORRADOR
	Limit         int     `json:"limit"`
}

// Sub-estructuras de respuesta
type AdminReportSummary struct {
	TotalEventos            int64   `json:"totalEventos"`
	TotalPublicados         int64   `json:"totalPublicados"`
	TotalCancelados         int64   `json:"totalCancelados"`
	TotalBorradores         int64   `json:"totalBorradores"`
	EntradasVendidasTotales int64   `json:"entradasVendidasTotales"`
	RecaudacionTotal        float64 `json:"recaudacionTotal"`
}

type AdminReportEvent struct {
	IdEvento         int64   `json:"idEvento" gorm:"column:evento_id"`
	Titulo           string  `json:"titulo"`
	Categoria        string  `json:"categoria"`
	Lugar            string  `json:"lugar"`
	Estado           string  `json:"estado"`
	FechaInicio      string  `json:"fechaInicio"`
	FechaFin         string  `json:"fechaFin"`
	EntradasVendidas int64   `json:"entradasVendidas"`
	RecaudacionTotal float64 `json:"recaudacionTotal"`
}

type AdminReportTopEvent struct {
	IdEvento         int64   `json:"idEvento"`
	Titulo           string  `json:"titulo"`
	Lugar            string  `json:"lugar"`
	EntradasVendidas int64   `json:"entradasVendidas"`
	Recaudacion      float64 `json:"recaudacion"`
}

type AdminReportCategory struct {
	IdCategoria      int64   `json:"idCategoria"`
	Categoria        string  `json:"categoria"`
	CantidadEventos  int64   `json:"cantidadEventos"`
	RecaudacionTotal float64 `json:"recaudacionTotal"`
	EntradasVendidas int64   `json:"entradasVendidas"`
}

// AdminReportResponse: Respuesta final consolidada
type AdminReportResponse struct {
	Summary    AdminReportSummary    `json:"summary"`
	Events     []AdminReportEvent    `json:"events"`
	TopEventos []AdminReportTopEvent `json:"topEventos"`
	ByCategory []AdminReportCategory `json:"byCategory"`
}
