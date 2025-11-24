package schemas

type EditarFechaEventoRequest struct {
	// ID del registro evento_fecha (evento_fecha_id)
	IdFechaEvento int64 `json:"idFechaEvento"`

	// Si quieres cambiar la fecha del calendario (tabla FECHA)
	// Necesitas el id de la fila fecha_id y la nueva fecha
	IdFecha    *int64  `json:"idFecha,omitempty"`    // opcional
	NuevaFecha *string `json:"nuevaFecha,omitempty"` // "YYYY-MM-DD"

	// Si quieres cambiar la HORA de inicio de esa función
	NuevaHoraInicio *string `json:"nuevaHoraInicio,omitempty"` // "HH:MM"

	// Si quieres reasignar este evento_fecha a otra fecha_id
	NuevoFechaID *int64 `json:"nuevoFechaId,omitempty"`
}

type EditarSectorRequest struct {
	IdSector      int64   `json:"idSector"`
	SectorTipo    *string `json:"sectorTipo,omitempty"`
	TotalEntradas *int    `json:"totalEntradas,omitempty"`
	CantVendidas  *int    `json:"cantVendidas,omitempty"` // solo uso admin
	Estado        *int16  `json:"estado,omitempty"`
}

type EditarPerfilRequest struct {
	IdPerfil int64   `json:"idPerfil"`
	Nombre   *string `json:"nombre,omitempty"`
	Estado   *int16  `json:"estado,omitempty"`
}

type EditarTipoTicketRequest struct {
	IdTipoTicket int64   `json:"idTipoTicket"`
	Nombre       *string `json:"nombre,omitempty"`
	FechaIni     *string `json:"fechaIni,omitempty"` // "YYYY-MM-DD"
	FechaFin     *string `json:"fechaFin,omitempty"` // "YYYY-MM-DD"
	Estado       *int16  `json:"estado,omitempty"`
}

type EditarEventoRequest struct {
	IdEvento int64 `json:"idEvento"`

	// Nivel EVENTO
	NuevoLugar          *string `json:"nuevoLugar,omitempty"`
	NuevoEstadoWorkflow *int16  `json:"nuevoEstadoWorkflow,omitempty"` // evento_estado
	NuevoEstadoFlag     *int16  `json:"nuevoEstadoFlag,omitempty"`     // estado (on/off)

	// Niveles relacionados:
	Fechas      []EditarFechaEventoRequest `json:"fechas,omitempty"`
	Sectores    []EditarSectorRequest      `json:"sectores,omitempty"`
	Perfiles    []EditarPerfilRequest      `json:"perfiles,omitempty"`
	TiposTicket []EditarTipoTicketRequest  `json:"tiposTicket,omitempty"`

	// Auditoría
	UsuarioModificacion int64 `json:"usuarioModificacion"`
}
