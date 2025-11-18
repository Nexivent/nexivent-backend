package adapter

import (
	goerrors "errors"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	util "github.com/Nexivent/nexivent-backend/internal/dao/model/util"
	daoPostgresql "github.com/Nexivent/nexivent-backend/internal/dao/repository"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
	"github.com/jackc/pgx/v5/pgconn"
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
		FechaInicio:     cuponReq.FechaInicio,
		FechaFin:        cuponReq.FechaFin,
		UsuarioCreacion: &usuario.ID,
		EventoID:        cuponReq.EventoID,
	}

	result := c.DaoPostgresql.Cupon.CrearCupon(cuponModel)

	if result != nil {
		// Intentamos convertir el error a un PgError (Propio de Postgres)
		var pgErr *pgconn.PgError

		if goerrors.As(result, &pgErr) {
			switch pgErr.Code {
			// Violación de UNIQUE → cupón ya existe
			case "23505":
				return nil, &errors.ConflictError.CuponAlreadyExists
			// Violación de Foreign Key
			case "23503":
				return nil, &errors.UnprocessableEntityError.InvalidEventoId // o el FK que corresponda

			// Violación de CHECK o restricciones de dominio
			case "23514":
				return nil, &errors.UnprocessableEntityError.InvalidRequestBody
			}
		}

		// Otros errores no controlados → error 500
		return nil, &errors.InternalServerError.Default
	}

	cuponRes := &schemas.CuponResponse{
		ID:            cuponModel.ID,
		Descripcion:   cuponModel.Descripcion,
		Tipo:          util.TipoCupon(cuponModel.Tipo),
		Valor:         cuponModel.Valor,
		Codigo:        cuponModel.Codigo,
		UsoPorUsuario: cuponReq.UsoPorUsuario,
		FechaInicio:   cuponReq.FechaInicio,
		FechaFin:      cuponReq.FechaFin,
	}

	return cuponRes, nil
}
