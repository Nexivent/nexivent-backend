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

type EventRequest struct {
	IdOrganizador     int     `json:"idOrganizador"`
	IdCategoria       int     `json:"idCategoria"`
	Titulo            string  `json:"titulo"`
	Descripcion       string  `json:"descripcion"`
	Lugar             string  `json:"lugar"`
	Estado            string  `json:"estado"`
	Likes             int     `json:"likes"`
	NoInteres         int     `json:"noInteres"`
	CantVendidasTotal int     `json:"cantVendidasTotal"`
	TotalRecaudado    float64 `json:"totalRecaudado"`
	ImagenPortada     string  `json:"imagenPortada"`
	ImagenLugar       string  `json:"imagenLugar"`
	VideoUrl          string  `json:"videoUrl"`

	Fechas      []EventoFecha     `json:"eventDates"`
	Perfiles    []PerfilDePersona `json:"perfiles"`
	Sectores    []Sector          `json:"sectores"`
	TiposTicket []TipoDeTicket    `json:"tiposTicket"`
	//Precios     map[string]PreciosSector `json:"precios"`

	// RELACIÓN 1–N: un evento tiene muchos cupones
	Cupones []Cupon `gorm:"foreignKey:EventoID;references:ID"`
}
