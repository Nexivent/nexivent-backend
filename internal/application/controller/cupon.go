package controller

import (
	"time"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/application/adapter"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
)

type CuponController struct {
	Logger       logging.Logger
	CuponAdapter *adapter.Cupon
}

func NewCuponController(
	logger logging.Logger,
	cuponAdapter *adapter.Cupon,
) *CuponController {
	return &CuponController{
		Logger:       logger,
		CuponAdapter: cuponAdapter,
	}
}

func (cc *CuponController) CreateCupon(
	cuponReq schemas.CuponResquest,
	usuarioCreacion int64,
) (*schemas.CuponResponse, *errors.Error) {
	return cc.CuponAdapter.CreatePostgresqlCupon(&cuponReq, usuarioCreacion)
}

func (cc *CuponController) UpdateCupon(
	cuponReq schemas.CuponResquest,
	usuarioModificacion int64,
) (*schemas.CuponResponse, *errors.Error) {
	return cc.CuponAdapter.UpdatePostgresqlCupon(&cuponReq, usuarioModificacion)
}

func (cc *CuponController) FetchCuponPorOrganizador(organizadorId int64) (*schemas.CuponesOrganizator, *errors.Error) {
	return cc.CuponAdapter.FetchPostresqlCuponPorOrganizador(organizadorId)
}

func (cc *CuponController) FetchValidarCuponParaOrdenDeCompra(usuarioId int64, fechaActual time.Time, eventoId int64, codigoCupon string) (*schemas.CuponResponseOrdenDePago, *errors.Error) {
	return cc.CuponAdapter.FetchPostresqlValidarCuponParaOrdenDeCompra(usuarioId, fechaActual, eventoId, codigoCupon)
}

func (cc *CuponController) CreateUsuarioCuponParaOrdenDeCompra(usarioCuponReq *schemas.UsuarioCuponRes) (*schemas.UsuarioCuponRes, *errors.Error) {
	return cc.CuponAdapter.CreatePostgresqlUsuarioCuponParaOrdenDeCompra(usarioCuponReq)
}
