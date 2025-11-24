package adapter

import (
	goerrors "errors"
	"time"

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

func (c *Cupon) UpdatePostgresqlCupon(cuponReq *schemas.CuponResquest, usuarioModificacion int64) (*schemas.CuponResponse, *errors.Error) {
	usuario, error := c.DaoPostgresql.Usuario.ObtenerUsuarioBasicoPorID(usuarioModificacion)

	if error != nil {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	_, errorCupon := c.DaoPostgresql.Cupon.ObtenerCuponPorIdYIdEvento(cuponReq.ID, cuponReq.EventoID)

	if errorCupon != nil {
		return nil, &errors.ObjectNotFoundError.CuponNotFound
	}

	now := time.Now()

	cuponModel := &model.Cupon{
		ID:                  cuponReq.ID,
		Descripcion:         cuponReq.Descripcion,
		Tipo:                cuponReq.Tipo.Codigo(),
		Valor:               cuponReq.Valor,
		EstadoCupon:         cuponReq.EstadoCupon.Codigo(), //activo
		Codigo:              cuponReq.Codigo,
		UsoPorUsuario:       cuponReq.UsoPorUsuario,
		FechaInicio:         cuponReq.FechaInicio,
		FechaFin:            cuponReq.FechaFin,
		UsuarioModificacion: &usuario.ID,
		FechaModificacion:   &now,
		EventoID:            cuponReq.EventoID,
	}

	result := c.DaoPostgresql.Cupon.ActualizarCupon(cuponModel)

	if result != nil {

		// Detectar si es error de Postgres
		var pgErr *pgconn.PgError
		if goerrors.As(result, &pgErr) {
			switch pgErr.Code {
			case "23505": // UNIQUE violation
				return nil, &errors.ConflictError.CuponAlreadyExists
			case "23514": // CHECK violation
				return nil, &errors.UnprocessableEntityError.InvalidRequestBody
			}
		}

		// Otros errores no controlados
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

func (c *Cupon) FetchPostresqlCuponPorOrganizador(oranizadorId int64) (*schemas.CuponesOrganizator, *errors.Error) {
	_, error := c.DaoPostgresql.Usuario.ObtenerUsuarioBasicoPorID(oranizadorId)

	if error != nil {
		return nil, &errors.ObjectNotFoundError.UserNotFound
	}

	cupones, result := c.DaoPostgresql.Cupon.ObtenerCuponesPorOrganizador(oranizadorId)

	if result != nil {
		return nil, &errors.InternalServerError.Default
	}

	var listCupones []*schemas.CuponOrganizator

	for _, cu := range cupones {
		listCupones = append(listCupones, &schemas.CuponOrganizator{
			ID:            cu.ID,
			Descripcion:   cu.Descripcion,
			Tipo:          util.TipoCupon(cu.Tipo),
			EstadoCupon:   util.Estado(cu.EstadoCupon),
			Valor:         cu.Valor,
			Codigo:        cu.Codigo,
			UsoPorUsuario: cu.UsoPorUsuario,
			UsoRealizados: cu.UsoRealizados,
			FechaInicio:   cu.FechaInicio,
			FechaFin:      cu.FechaFin,
			EventoID:      cu.EventoID,
		})
	}

	cuponesRes := &schemas.CuponesOrganizator{
		Cupones: listCupones,
	}

	return cuponesRes, nil
}

func (c *Cupon) FetchPostresqlValidarCuponParaOrdenDeCompra(usuarioId int64, fechaActual time.Time, eventoId int64, codigoCupon string) (*schemas.CuponResponseOrdenDePago, *errors.Error) {
	_, usuarioErr := c.DaoPostgresql.Usuario.ObtenerUsuarioBasicoPorID(usuarioId)

	if usuarioErr != nil { // el usuario no existe
		return nil, &errors.ObjectNotFoundError.UserNotFound
	}

	_, eventoErr := c.DaoPostgresql.Evento.ObtenerEventoPorId(eventoId)

	if eventoErr != nil { // el evento no existe
		return nil, &errors.ObjectNotFoundError.EventoNotFound
	}

	cuponModel, cuponEventoErr := c.DaoPostgresql.Cupon.ObtenerCuponPorCodYIdEvento(eventoId, codigoCupon)

	if cuponEventoErr != nil { // el cupón no existe o no existe en este evento
		return nil, &errors.ObjectNotFoundError.CuponNotFound
	}

	if cuponModel == nil {
		c.logger.Error("cuponModel es nil después de búsqueda sin error", "eventoId", eventoId, "codigo", codigoCupon)
		return nil, &errors.ObjectNotFoundError.CuponNotFound
	}

	if cuponModel.FechaInicio.IsZero() || cuponModel.FechaFin.IsZero() {
		c.logger.Error("Cupón tiene fechas inválidas", "cuponId", cuponModel.ID)
		return nil, &errors.BadRequestError.InvalidFechaCupon
	}

	if fechaActual.After(cuponModel.FechaFin) || fechaActual.Before(cuponModel.FechaInicio) {
		//No se puede usar el cupón en esta fecha
		return nil, &errors.BadRequestError.InvalidFechaCupon
	}

	usuarioCuponModel, usuarioCuponErr := c.DaoPostgresql.UsuarioCupon.ObtenerUsuarioCuponPorId(usuarioId, cuponModel.ID)

	cuponRes := &schemas.CuponResponseOrdenDePago{
		ID:        cuponModel.ID,
		Tipo:      util.TipoCupon(cuponModel.Tipo),
		Valor:     cuponModel.Valor,
		CantUsada: 0,
	}

	if usuarioCuponErr != nil {
		//el usuario nunca usó el cupón
		return cuponRes, nil
	}

	cuponRes.CantUsada = usuarioCuponModel.CantUsada

	if usuarioCuponModel.CantUsada >= cuponModel.UsoPorUsuario {
		// el usuario ya utilizó el cupón hasta el límite de veces válidas
		return nil, &errors.BadRequestError.CantLimitUseCupon
	}

	return cuponRes, nil
}

func (c *Cupon) CreatePostgresqlUsuarioCuponParaOrdenDeCompra(usarioCuponReq *schemas.UsuarioCuponRes) (*schemas.UsuarioCuponRes, *errors.Error) {

	usuarioCuponModel := &model.UsuarioCupon{
		CuponID:   usarioCuponReq.UsuarioID,
		UsuarioID: usarioCuponReq.UsuarioID,
		CantUsada: usarioCuponReq.CantUsada,
	}

	usuarioCuponErr := c.DaoPostgresql.UsuarioCupon.ActualizarUsuarioCupon(usuarioCuponModel)

	if usuarioCuponErr != nil { // no existe
		return nil, &errors.BadRequestError.UsuarioCuponNotUpdate
	}

	return usarioCuponReq, nil
}
