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

type PerfilPersonaAdapter struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.NexiventPsqlEntidades
}

func NewPerfilPersonaAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.NexiventPsqlEntidades,
) *PerfilPersonaAdapter {
	return &PerfilPersonaAdapter{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

func (a *PerfilPersonaAdapter) CrearPerfilPersona(req *schemas.PerfilPersonaRequest, usuarioCreacion int64) (*schemas.PerfilPersonaResponse, *errors.Error) {
	now := time.Now()

	modelo := &model.PerfilDePersona{
		EventoID:        req.EventoID,
		Nombre:          req.Nombre,
		Estado:          req.Estado,
		UsuarioCreacion: &usuarioCreacion,
		FechaCreacion:   now,
	}

	if err := a.DaoPostgresql.PerfilDePersona.CrearPerfilDePersona(modelo); err != nil {
		a.logger.Errorf("CrearPerfilPersona: %v", err)
		return nil, &errors.BadRequestError.EventoNotCreated // puedes crear uno específico si quieres
	}

	resp := &schemas.PerfilPersonaResponse{
		ID:       modelo.ID,
		EventoID: modelo.EventoID,
		Nombre:   modelo.Nombre,
		Estado:   modelo.Estado,
	}
	return resp, nil
}

func (a *PerfilPersonaAdapter) ActualizarPerfilPersona(id int64, req *schemas.PerfilPersonaUpdateRequest, usuarioModificacion int64) (*schemas.PerfilPersonaResponse, *errors.Error) {
	now := time.Now()

	perfil, err := a.DaoPostgresql.PerfilDePersona.ModificarPerfilDePersonaPorCampos(
		id,
		req.EventoID,
		req.Nombre,
		req.Estado,
		&usuarioModificacion,
		&now,
	)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}
		a.logger.Errorf("ActualizarPerfilPersona(%d): %v", id, err)
		return nil, &errors.BadRequestError.EventoNotUpdated
	}

	resp := &schemas.PerfilPersonaResponse{
		ID:       perfil.ID,
		EventoID: perfil.EventoID,
		Nombre:   perfil.Nombre,
		Estado:   perfil.Estado,
	}
	return resp, nil
}

func (a *PerfilPersonaAdapter) ListarPerfilesPorEvento(eventoID int64) ([]schemas.PerfilPersonaResponse, *errors.Error) {
	perfiles, err := a.DaoPostgresql.PerfilDePersona.ListarPerfilPersonaPorEventoID(eventoID)
	if err != nil {
		a.logger.Errorf("ListarPerfilesPorEvento(%d): %v", eventoID, err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	out := make([]schemas.PerfilPersonaResponse, len(perfiles))
	for i, p := range perfiles {
		out[i] = schemas.PerfilPersonaResponse{
			ID:       p.ID,
			EventoID: p.EventoID,
			Nombre:   p.Nombre,
			Estado:   p.Estado,
		}
	}
	return out, nil
}
