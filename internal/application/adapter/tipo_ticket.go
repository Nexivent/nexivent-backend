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

type TipoTicketAdapter struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.NexiventPsqlEntidades
}

func NewTipoTicketAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.NexiventPsqlEntidades,
) *TipoTicketAdapter {
	return &TipoTicketAdapter{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

func (a *TipoTicketAdapter) CrearTipoTicket(req *schemas.TipoTicketTicketRequest, usuarioCreacion int64) (*schemas.TipoTicketTicketResponse, *errors.Error) {
	now := time.Now()

	fechaIni, err := time.Parse("2006-01-02", req.FechaIni)
	if err != nil {
		return nil, &errors.UnprocessableEntityError.InvalidDateFormat
	}
	fechaFin, err := time.Parse("2006-01-02", req.FechaFin)
	if err != nil {
		return nil, &errors.UnprocessableEntityError.InvalidDateFormat
	}

	modelo := &model.TipoDeTicket{
		EventoID:        req.EventoID,
		Nombre:          req.Nombre,
		FechaIni:        fechaIni,
		FechaFin:        fechaFin,
		Estado:          req.Estado,
		UsuarioCreacion: &usuarioCreacion,
		FechaCreacion:   now,
	}

	if err := a.DaoPostgresql.TipoDeTicket.CrearTipoDeTicket(modelo); err != nil {
		a.logger.Errorf("CrearTipoTicket: %v", err)
		return nil, &errors.BadRequestError.EventoNotCreated
	}

	resp := &schemas.TipoTicketTicketResponse{
		ID:       modelo.ID,
		EventoID: modelo.EventoID,
		Nombre:   modelo.Nombre,
		FechaIni: modelo.FechaIni.Format("2006-01-02"),
		FechaFin: modelo.FechaFin.Format("2006-01-02"),
		Estado:   modelo.Estado,
	}
	return resp, nil
}

func (a *TipoTicketAdapter) ActualizarTipoTicket(id int64, req *schemas.TipoTicketUpdateRequest, usuarioModificacion int64) (*schemas.TipoTicketTicketResponse, *errors.Error) {
	now := time.Now()

	var fechaIni *time.Time
	var fechaFin *time.Time

	if req.FechaIni != nil {
		parsed, err := time.Parse("2006-01-02", *req.FechaIni)
		if err != nil {
			return nil, &errors.UnprocessableEntityError.InvalidDateFormat
		}
		fechaIni = &parsed
	}
	if req.FechaFin != nil {
		parsed, err := time.Parse("2006-01-02", *req.FechaFin)
		if err != nil {
			return nil, &errors.UnprocessableEntityError.InvalidDateFormat
		}
		fechaFin = &parsed
	}

	// No tienes ModificarTipoTicketPorCampos, así que lo hago cargando y guardando:
	var modelo model.TipoDeTicket
	db := a.DaoPostgresql.TipoDeTicket.PostgresqlDB
	if err := db.First(&modelo, "tipo_de_ticket_id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}
		a.logger.Errorf("ActualizarTipoTicket First(%d): %v", id, err)
		return nil, &errors.BadRequestError.EventoNotUpdated
	}

	if req.Nombre != nil {
		modelo.Nombre = *req.Nombre
	}
	if fechaIni != nil {
		modelo.FechaIni = *fechaIni
	}
	if fechaFin != nil {
		modelo.FechaFin = *fechaFin
	}
	if req.Estado != nil {
		modelo.Estado = *req.Estado
	}
	modelo.UsuarioModificacion = &usuarioModificacion
	modelo.FechaModificacion = &now

	if err := db.Save(&modelo).Error; err != nil {
		a.logger.Errorf("ActualizarTipoTicket Save(%d): %v", id, err)
		return nil, &errors.BadRequestError.EventoNotUpdated
	}

	resp := &schemas.TipoTicketTicketResponse{
		ID:       modelo.ID,
		EventoID: modelo.EventoID,
		Nombre:   modelo.Nombre,
		FechaIni: modelo.FechaIni.Format("2006-01-02"),
		FechaFin: modelo.FechaFin.Format("2006-01-02"),
		Estado:   modelo.Estado,
	}
	return resp, nil
}

func (a *TipoTicketAdapter) ListarTiposTicketPorEvento(eventoID int64) ([]schemas.TipoTicketTicketResponse, *errors.Error) {
	// Llamada al DAO (Recuerda el .Preload("Tarifa"))
	list, err := a.DaoPostgresql.TipoDeTicket.ListarTipoTicketPorEventoID(eventoID)
	if err != nil {
		a.logger.Errorf("ListarTiposTicketPorEvento(%d): %v", eventoID, err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	out := make([]schemas.TipoTicketTicketResponse, len(list))
	for i, t := range list {

		// --- INICIO: Mapeo de Tarifas ---
		var tarifasResp []schemas.TarifaResponseOtros
		if len(t.Tarifa) > 0 {
			tarifasResp = make([]schemas.TarifaResponseOtros, len(t.Tarifa))
			for j, tar := range t.Tarifa {
				tarifasResp[j] = schemas.TarifaResponseOtros{
					ID:     tar.ID,
					Precio: tar.Precio, // Ajusta según los campos reales
					Estado: tar.Estado, // Ejemplo
					// Mapea aquí el resto de campos de tarifa
				}
			}
		}
		// --- FIN: Mapeo de Tarifas ---

		out[i] = schemas.TipoTicketTicketResponse{
			ID:       t.ID,
			EventoID: t.EventoID,
			Nombre:   t.Nombre,
			FechaIni: t.FechaIni.Format("2006-01-02"),
			FechaFin: t.FechaFin.Format("2006-01-02"),
			Estado:   t.Estado,
			Tarifas:  tarifasResp, // Asignamos el array mapeado
		}
	}
	return out, nil
}
