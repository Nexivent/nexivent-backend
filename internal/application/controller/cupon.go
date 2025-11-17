package controller

import (
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
