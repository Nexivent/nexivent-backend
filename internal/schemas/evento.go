package schemas

import "github.com/Loui27/nexivent-backend/internal/dao/model"

type EventosPaginados struct {
	Eventos      []*model.Evento `json:"eventos"`
	Total        int64           `json:"total"`
	PaginaActual int             `json:"pagina_actual"`
	TotalPaginas int             `json:"total_paginas"`
}
