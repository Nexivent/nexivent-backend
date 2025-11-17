package schemas

import (
	util "github.com/Loui27/nexivent-backend/internal/dao/model/util"
)

type CuponResponse struct { //response // API -> front
	ID            int64
	Descripcion   string
	Tipo          util.TipoCupon
	Valor         float64
	Codigo        string
	UsoPorUsuario int64
	//EventoID      int64
}

type CuponResquest struct { //request //front -> API
	ID            int64
	Descripcion   string
	Tipo          util.TipoCupon
	Valor         float64
	Codigo        string
	UsoPorUsuario int64
	EventoID      int64
}
