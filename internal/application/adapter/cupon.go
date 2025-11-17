package adapter

import (
	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	util "github.com/Nexivent/nexivent-backend/internal/dao/model/util"
	daoPostgresql "github.com/Nexivent/nexivent-backend/internal/dao/repository"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
)

type Cupon struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.NexiventPsqlEntidades
}

func NewCuponAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.NexiventPsqlEntidades,
) *Cupon {
	return &Cupon{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

func (c *Cupon) CreatePostgresqlCupon(cuponReq *schemas.CuponResquest, usuarioCreacion int64) (*schemas.CuponResponse, *errors.Error) {
	usuario, error := c.DaoPostgresql.Usuario.ObtenerUsuarioBasicoPorID(usuarioCreacion)

	if error != nil {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	cuponModel := &model.Cupon{
		Descripcion:     cuponReq.Descripcion,
		Tipo:            cuponReq.Tipo.Codigo(),
		Valor:           cuponReq.Valor,
		EstadoCupon:     util.Activo.Codigo(), //activo
		Codigo:          cuponReq.Codigo,
		UsoPorUsuario:   cuponReq.UsoPorUsuario,
		UsoRealizados:   0, // sin uso aún
		UsuarioCreacion: &usuario.ID,
	}

	result := c.DaoPostgresql.Cupon.CrearCupon(cuponModel)
	if result != nil {
		return nil, &errors.ConflictError.CuponAlreadyExits
	}

	cuponRes := &schemas.CuponResponse{
		ID:            cuponModel.ID,
		Descripcion:   cuponModel.Descripcion,
		Tipo:          util.TipoCupon(cuponModel.Tipo),
		Valor:         cuponModel.Valor,
		Codigo:        cuponModel.Codigo,
		UsoPorUsuario: cuponReq.UsoPorUsuario,
	}

	return cuponRes, nil
}
