package adapter

import (
	"time"

	"github.com/Nexivent/nexivent-backend/errors"
	model "github.com/Nexivent/nexivent-backend/internal/dao/model"
	daoPostgresql "github.com/Nexivent/nexivent-backend/internal/dao/repository"
	schemas "github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
)

type SectorAdapter struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.NexiventPsqlEntidades
}

func NewSectorAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.NexiventPsqlEntidades,
) *SectorAdapter {
	return &SectorAdapter{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

func (a *SectorAdapter) CrearSector(req *schemas.SectorTicketRequest, usuarioCreacion int64) (*schemas.SectorTicketResponse, *errors.Error) {
	now := time.Now()

	modelo := &model.Sector{
		EventoID:        req.EventoID,
		SectorTipo:      req.SectorTipo,
		TotalEntradas:   req.TotalEntradas,
		CantVendidas:    0,
		Estado:          req.Estado,
		UsuarioCreacion: &usuarioCreacion,
		FechaCreacion:   now,
	}

	if err := a.DaoPostgresql.Sector.CrearSector(modelo); err != nil {
		a.logger.Errorf("CrearSector: %v", err)
		return nil, &errors.BadRequestError.EventoNotCreated
	}

	resp := &schemas.SectorTicketResponse{
		ID:            modelo.ID,
		EventoID:      modelo.EventoID,
		SectorTipo:    modelo.SectorTipo,
		TotalEntradas: modelo.TotalEntradas,
		CantVendidas:  modelo.CantVendidas,
		Estado:        modelo.Estado,
	}
	return resp, nil
}

func (a *SectorAdapter) ActualizarSector(id int64, req *schemas.SectorUpdateRequest, usuarioModificacion int64) (*schemas.SectorTicketResponse, *errors.Error) {
	now := time.Now()

	sector, err := a.DaoPostgresql.Sector.ModificarSectorPorCampos(
		id,
		req.SectorTipo,
		req.TotalEntradas,
		req.CantVendidas,
		req.Estado,
		&usuarioModificacion,
		&now,
	)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}
		a.logger.Errorf("ActualizarSector(%d): %v", id, err)
		return nil, &errors.BadRequestError.EventoNotUpdated
	}

	resp := &schemas.SectorTicketResponse{
		ID:            sector.ID,
		EventoID:      sector.EventoID,
		SectorTipo:    sector.SectorTipo,
		TotalEntradas: sector.TotalEntradas,
		CantVendidas:  sector.CantVendidas,
		Estado:        sector.Estado,
	}
	return resp, nil
}

func (a *SectorAdapter) ListarSectoresPorEvento(eventoID int64) ([]schemas.SectorTicketResponse, *errors.Error) {
	sectores, err := a.DaoPostgresql.Sector.ListarSectorePorIdEvento(eventoID)
	if err != nil {
		a.logger.Errorf("ListarSectoresPorEvento(%d): %v", eventoID, err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	out := make([]schemas.SectorTicketResponse, len(sectores))
	for i, s := range sectores {
		out[i] = schemas.SectorTicketResponse{
			ID:            s.ID,
			EventoID:      s.EventoID,
			SectorTipo:    s.SectorTipo,
			TotalEntradas: s.TotalEntradas,
			CantVendidas:  s.CantVendidas,
			Estado:        s.Estado,
		}
	}
	return out, nil
}
