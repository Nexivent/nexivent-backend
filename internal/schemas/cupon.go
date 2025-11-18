package schemas

import (
	"time"

	util "github.com/Nexivent/nexivent-backend/internal/dao/model/util"
)

type CuponResponse struct { //response // API -> front
	ID            int64
	Descripcion   string
	Tipo          util.TipoCupon
	Valor         float64
	Codigo        string
	UsoPorUsuario int64
	FechaInicio   time.Time
	FechaFin      time.Time
	//EventoID      int64
}

type CuponResquest struct { //request //front -> API
	ID            int64
	Descripcion   string
	Tipo          util.TipoCupon
	Valor         float64
	Codigo        string
	EstadoCupon   util.Estado
	UsoPorUsuario int64
	FechaInicio   time.Time
	FechaFin      time.Time
	EventoID      int64
}

type CuponOrganizator struct {
	ID            int64
	Descripcion   string
	Tipo          util.TipoCupon
	EstadoCupon   util.Estado
	Valor         float64
	Codigo        string
	UsoPorUsuario int64
	UsoRealizados int64
	FechaInicio   time.Time
	FechaFin      time.Time
	EventoID      int64
}

type CuponesOrganizator struct {
	Cupones []*CuponOrganizator
}
