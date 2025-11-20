package model

import "time"

type ValidacionDocumento struct {
    ID               int64     `json:"id" db:"id"`
    TipoDocumento    string    `json:"tipo_documento" db:"tipo_documento"` // DNI, RUC
    NumeroDocumento  string    `json:"numero_documento" db:"numero_documento"`
    NombreCompleto   string    `json:"nombre_completo" db:"nombre_completo"`
    RazonSocial      string    `json:"razon_social,omitempty" db:"razon_social"`
    Direccion        string    `json:"direccion,omitempty" db:"direccion"`
    Departamento     string    `json:"departamento,omitempty" db:"departamento"`
    Provincia        string    `json:"provincia,omitempty" db:"provincia"`
    Distrito         string    `json:"distrito,omitempty" db:"distrito"`
    Ubigeo           string    `json:"ubigeo,omitempty" db:"ubigeo"`
    EstadoContribuyente string `json:"estado_contribuyente,omitempty" db:"estado_contribuyente"`
    CondicionContribuyente string `json:"condicion_contribuyente,omitempty" db:"condicion_contribuyente"`
    FechaConsulta    time.Time `json:"fecha_consulta" db:"fecha_consulta"`
    Valido           bool      `json:"valido" db:"valido"`
}