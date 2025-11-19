package schemas

type ValidarDocumentoRequest struct {
	TipoDocumento   string `json:"tipo_documento"` // DNI, CE, RUC_PERSONA, RUC_EMPRESA
	NumeroDocumento string `json:"numero_documento"`
}

type ValidarDocumentoResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    *ValidacionDocumentoData `json:"data,omitempty"`
}

type ValidacionDocumentoData struct {
	TipoDocumento          string `json:"tipo_documento"`
	NumeroDocumento        string `json:"numero_documento"`
	NombreCompleto         string `json:"nombre_completo,omitempty"`
	RazonSocial            string `json:"razon_social,omitempty"`
	Direccion              string `json:"direccion,omitempty"`
	Departamento           string `json:"departamento,omitempty"`
	Provincia              string `json:"provincia,omitempty"`
	Distrito               string `json:"distrito,omitempty"`
	Ubigeo                 string `json:"ubigeo,omitempty"`
	EstadoContribuyente    string `json:"estado_contribuyente,omitempty"`
	CondicionContribuyente string `json:"condicion_contribuyente,omitempty"`
	EsEmpresa              bool   `json:"es_empresa"`
	Valido                 bool   `json:"valido"`
}