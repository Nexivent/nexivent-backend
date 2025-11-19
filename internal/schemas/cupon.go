package schemas

import (
	"time"

	util "github.com/Nexivent/nexivent-backend/internal/dao/model/util"
)

// response // API -> front
type CuponResponse struct {
	ID            int64          `json:"id"`
	Descripcion   string         `json:"descripcion"`
	Tipo          util.TipoCupon `json:"tipo"`
	Valor         float64        `json:"valor"`
	Codigo        string         `json:"codigo"`
	UsoPorUsuario int64          `json:"usoPorUsuario"`
	FechaInicio   time.Time      `json:"fechaInicio"`
	FechaFin      time.Time      `json:"fechaFin"`
	//EventoID    int64         `json:"eventoId,omitempty"`
}

// request //front -> API
type CuponResquest struct {
	ID            int64          `json:"id"`
	Descripcion   string         `json:"descripcion"`
	Tipo          util.TipoCupon `json:"tipo"`
	Valor         float64        `json:"valor"`
	Codigo        string         `json:"codigo"`
	EstadoCupon   util.Estado    `json:"estadoCupon"`
	UsoPorUsuario int64          `json:"usoPorUsuario"`
	FechaInicio   time.Time      `json:"fechaInicio"`
	FechaFin      time.Time      `json:"fechaFin"`
	EventoID      int64          `json:"eventoId"`
}

type CuponOrganizator struct {
	ID            int64          `json:"id"`
	Descripcion   string         `json:"descripcion"`
	Tipo          util.TipoCupon `json:"tipo"`
	EstadoCupon   util.Estado    `json:"estadoCupon"`
	Valor         float64        `json:"valor"`
	Codigo        string         `json:"codigo"`
	UsoPorUsuario int64          `json:"usoPorUsuario"`
	UsoRealizados int64          `json:"usoRealizados"`
	FechaInicio   time.Time      `json:"fechaInicio"`
	FechaFin      time.Time      `json:"fechaFin"`
	EventoID      int64          `json:"eventoId"`
}

type CuponesOrganizator struct {
	Cupones []*CuponOrganizator `json:"cupones"`
}
