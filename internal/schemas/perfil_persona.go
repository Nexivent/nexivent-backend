package schemas

type PerfilPersonaRequest struct {
	EventoID int64  `json:"idEvento"`
	Nombre   string `json:"nombre"`
	Estado   int16  `json:"estado"` // normalmente 1
}

type PerfilPersonaUpdateRequest struct {
	EventoID *int64  `json:"idEvento,omitempty"`
	Nombre   *string `json:"nombre,omitempty"`
	Estado   *int16  `json:"estado,omitempty"`
}

type PerfilPersonaResponse struct {
	ID       int64  `json:"idPerfilPersona"`
	EventoID int64  `json:"idEvento"`
	Nombre   string `json:"nombre"`
	Estado   int16  `json:"estado"`
}
