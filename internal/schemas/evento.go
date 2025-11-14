package schemas

import (
	"github.com/Loui27/nexivent-backend/internal/dao/model"
)

type EventosPaginados struct {
	Eventos      []*model.Evento `json:"eventos"`
	Total        int64           `json:"total"`
	PaginaActual int             `json:"pagina_actual"`
	TotalPaginas int             `json:"total_paginas"`
}

// EventDateRequest represents the event date information in the request
type EventDateRequest struct {
	IdFechaEvento int64  `json:"idFechaEvento,omitempty"`
	IdFecha       int64  `json:"idFecha,omitempty"`
	Fecha         string `json:"fecha"`
	HoraInicio    string `json:"horaInicio"`
	HoraFin       string `json:"horaFin"`
}

// EventDateResponse represents the event date information in the response
type EventDateResponse struct {
	IdFechaEvento int64  `json:"idFechaEvento"`
	IdFecha       int64  `json:"idFecha"`
	Fecha         string `json:"fecha"`
	HoraInicio    string `json:"horaInicio"`
	HoraFin       string `json:"horaFin"`
}

// PerfilRequest represents profile information
type PerfilRequest struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// PerfilResponse represents profile information in the response
type PerfilResponse struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// SectorRequest represents sector information
type SectorRequest struct {
	ID        string `json:"id"`
	Nombre    string `json:"nombre"`
	Capacidad int    `json:"capacidad"`
}

// SectorResponse represents sector information in the response
type SectorResponse struct {
	ID        string `json:"id"`
	Nombre    string `json:"nombre"`
	Capacidad int    `json:"capacidad"`
}

// TipoTicketRequest represents ticket type information
type TipoTicketRequest struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// TipoTicketResponse represents ticket type information in the response
type TipoTicketResponse struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// PrecioDetalle represents the price for a specific ticket type
type PrecioDetalle map[string]float64

// PreciosPerfil represents prices for all ticket types for a specific profile
type PreciosPerfil map[string]PrecioDetalle

// PreciosSector represents prices for all profiles in a sector
type PreciosSector map[string]PreciosPerfil

// MetadataRequest represents metadata in the request
type MetadataRequest struct {
	CreadoPor           string `json:"creadoPor"`
	FechaCreacion       string `json:"fechaCreacion"`
	UltimaActualizacion string `json:"ultimaActualizacion"`
	Version             int    `json:"version"`
}

// MetadataResponse represents metadata in the response
type MetadataResponse struct {
	CreadoPor           string `json:"creadoPor"`
	FechaCreacion       string `json:"fechaCreacion"`
	UltimaActualizacion string `json:"ultimaActualizacion"`
	Version             int    `json:"version"`
}

// EventoRequest represents the request payload for creating/updating an event
type EventoRequest struct {
	IdOrganizador     int64             `json:"idOrganizador"`
	IdCategoria       int64             `json:"idCategoria"`
	Titulo            string            `json:"titulo"`
	Descripcion       string            `json:"descripcion"`
	Lugar             string            `json:"lugar"`
	Estado            string            `json:"estado"`
	Likes             int               `json:"likes"`
	NoInteres         int               `json:"noInteres"`
	CantVendidasTotal int               `json:"cantVendidasTotal"`
	TotalRecaudado    float64           `json:"totalRecaudado"`
	ImagenPortada     string            `json:"imagenPortada"`
	ImagenLugar       string            `json:"imagenLugar"`
	VideoUrl          string            `json:"videoUrl"`
	EventDates        []EventDateRequest `json:"eventDates"`
	Perfiles          []PerfilRequest   `json:"perfiles"`
	Sectores          []SectorRequest   `json:"sectores"`
	TiposTicket       []TipoTicketRequest `json:"tiposTicket"`
	Precios           PreciosSector     `json:"precios"`
	Metadata          MetadataRequest   `json:"metadata"`
}

// EventoResponse represents the response payload for an event
type EventoResponse struct {
	IdEvento          int64              `json:"idEvento"`
	IdOrganizador     int64              `json:"idOrganizador"`
	IdCategoria       int64              `json:"idCategoria"`
	Titulo            string             `json:"titulo"`
	Descripcion       string             `json:"descripcion"`
	Lugar             string             `json:"lugar"`
	Estado            string             `json:"estado"`
	Likes             int                `json:"likes"`
	NoInteres         int                `json:"noInteres"`
	CantVendidasTotal int                `json:"cantVendidasTotal"`
	TotalRecaudado    float64            `json:"totalRecaudado"`
	ImagenPortada     string             `json:"imagenPortada"`
	ImagenLugar       string             `json:"imagenLugar"`
	VideoUrl          string             `json:"videoUrl"`
	EventDates        []EventDateResponse `json:"eventDates"`
	Perfiles          []PerfilResponse   `json:"perfiles"`
	Sectores          []SectorResponse   `json:"sectores"`
	TiposTicket       []TipoTicketResponse `json:"tiposTicket"`
	Precios           PreciosSector      `json:"precios"`
	Metadata          MetadataResponse   `json:"metadata"`
}
