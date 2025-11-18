package adapter

import (
	"time"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	daoPostgresql "github.com/Nexivent/nexivent-backend/internal/dao/repository"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
)

type TarifaAdapter struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.NexiventPsqlEntidades
}

func NewTarifaAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.NexiventPsqlEntidades,
) *TarifaAdapter {
	return &TarifaAdapter{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

func (a *TarifaAdapter) CrearTarifa(req *schemas.TarifaRequest, usuarioCreacion int64) (*schemas.TarifaResponse, *errors.Error) {
	now := time.Now()

	modelo := &model.Tarifa{
		SectorID:          req.SectorID,
		TipoDeTicketID:    req.TipoDeTicketID,
		PerfilDePersonaID: req.PerfilDePersonaID,
		Precio:            req.Precio,
		Estado:            req.Estado,
		UsuarioCreacion:   &usuarioCreacion,
		FechaCreacion:     now,
	}

	if err := a.DaoPostgresql.Tarifa.CrearTarifa(modelo); err != nil {
		a.logger.Errorf("CrearTarifa: %v", err)
		return nil, &errors.BadRequestError.EventoNotCreated
	}

	resp := &schemas.TarifaResponse{
		ID:                modelo.ID,
		SectorID:          modelo.SectorID,
		TipoDeTicketID:    modelo.TipoDeTicketID,
		PerfilDePersonaID: modelo.PerfilDePersonaID,
		Precio:            modelo.Precio,
		Estado:            modelo.Estado,
	}
	return resp, nil
}

func (a *TarifaAdapter) ActualizarTarifa(id int64, req *schemas.TarifaUpdateRequest, usuarioModificacion int64) (*schemas.TarifaResponse, *errors.Error) {
	now := time.Now()

	tarifa, err := a.DaoPostgresql.Tarifa.ModificarTarifaPorCampos(
		id,
		req.SectorID,
		req.TipoDeTicketID,
		req.PerfilDePersonaID,
		req.Precio,
		req.Estado,
		&usuarioModificacion,
		&now,
	)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}
		a.logger.Errorf("ActualizarTarifa(%d): %v", id, err)
		return nil, &errors.BadRequestError.EventoNotUpdated
	}

	resp := &schemas.TarifaResponse{
		ID:                tarifa.ID,
		SectorID:          tarifa.SectorID,
		TipoDeTicketID:    tarifa.TipoDeTicketID,
		PerfilDePersonaID: tarifa.PerfilDePersonaID,
		Precio:            tarifa.Precio,
		Estado:            tarifa.Estado,
	}
	return resp, nil
}

// ListarTarifasPorIDs es útil para el front si necesita detalles de varias tarifas a la vez
func (a *TarifaAdapter) ListarTarifasPorIDs(ids []int64) ([]schemas.TarifaResponse, *errors.Error) {
	list, err := a.DaoPostgresql.Tarifa.ObtenerTarifasPorIDs(ids)
	if err != nil {
		a.logger.Errorf("ListarTarifasPorIDs: %v", err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	out := make([]schemas.TarifaResponse, len(list))
	for i, t := range list {
		out[i] = schemas.TarifaResponse{
			ID:                t.ID,
			SectorID:          t.SectorID,
			TipoDeTicketID:    t.TipoDeTicketID,
			PerfilDePersonaID: t.PerfilDePersonaID,
			Precio:            t.Precio,
			Estado:            t.Estado,
		}
	}
	return out, nil
}
