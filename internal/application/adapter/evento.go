package adapter

import (
	"fmt"
	"time"

	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	daoPostgresql "github.com/Nexivent/nexivent-backend/internal/dao/repository"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
	"github.com/Nexivent/nexivent-backend/utils/convert"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Evento struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.NexiventPsqlEntidades
}

// Creates Evento adapter
func NewEventoAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.NexiventPsqlEntidades,
) *Evento {
	return &Evento{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// CreatePostgresqlEvento creates a new event with all related entities
func (e *Evento) CreatePostgresqlEvento(eventoReq *schemas.EventoRequest, usuarioCreacion int64) (*schemas.EventoResponse, *errors.Error) {
	// Start a transaction
	tx := e.DaoPostgresql.Evento.PostgresqlDB.Begin()
	if tx.Error != nil {
		e.logger.Errorf("Failed to begin transaction: %v", tx.Error)
		return nil, &errors.BadRequestError.EventoNotCreated
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()

	// Create the main event model
	eventoModel := &model.Evento{
		OrganizadorID:     eventoReq.IdOrganizador,
		CategoriaID:       eventoReq.IdCategoria,
		Titulo:            eventoReq.Titulo,
		Descripcion:       eventoReq.Descripcion,
		Lugar:             eventoReq.Lugar,
		EventoEstado:      convert.MapEstadoToInt16(eventoReq.Estado),
		CantMeGusta:       eventoReq.Likes,
		CantNoInteresa:    eventoReq.NoInteres,
		CantVendidoTotal:  eventoReq.CantVendidasTotal,
		ImagenPortada:     eventoReq.ImagenPortada,
		ImagenEscenario:   eventoReq.ImagenLugar,
		VideoPresentacion: eventoReq.VideoUrl,
		TotalRecaudado:    eventoReq.TotalRecaudado,
		Estado:            1, // Active by default
		UsuarioCreacion:   &usuarioCreacion,
		FechaCreacion:     now,
	}

	// Create the event
	if err := tx.Create(eventoModel).Error; err != nil {
		tx.Rollback()
		e.logger.Errorf("Failed to create event: %v", err)
		return nil, &errors.BadRequestError.EventoNotCreated
	}

	// Create perfiles (profiles)
	perfilesMap := make(map[string]int64)
	for _, perfil := range eventoReq.Perfiles {
		perfilModel := &model.PerfilDePersona{
			EventoID:        eventoModel.ID,
			Nombre:          perfil.Label,
			Estado:          1,
			UsuarioCreacion: &usuarioCreacion,
			FechaCreacion:   now,
		}
		if err := tx.Create(perfilModel).Error; err != nil {
			tx.Rollback()
			e.logger.Errorf("Failed to create perfil: %v", err)
			return nil, &errors.BadRequestError.EventoNotCreated
		}
		perfilesMap[perfil.ID] = perfilModel.ID
	}

	// Create sectores (sectors)
	sectoresMap := make(map[string]int64)
	for _, sector := range eventoReq.Sectores {
		sectorModel := &model.Sector{
			EventoID:        eventoModel.ID,
			SectorTipo:      sector.Nombre,
			TotalEntradas:   sector.Capacidad,
			CantVendidas:    0,
			Estado:          1,
			UsuarioCreacion: &usuarioCreacion,
			FechaCreacion:   now,
		}
		if err := tx.Create(sectorModel).Error; err != nil {
			tx.Rollback()
			e.logger.Errorf("Failed to create sector: %v", err)
			return nil, &errors.BadRequestError.EventoNotCreated
		}
		sectoresMap[sector.ID] = sectorModel.ID
	}

	// Create tipos de ticket (ticket types)
	tiposTicketMap := make(map[string]int64)
	for _, tipoTicket := range eventoReq.TiposTicket {
		fechaIni := now
		fechaFin := now.AddDate(1, 0, 0) // Default: 1 year from now

		tipoTicketModel := &model.TipoDeTicket{
			EventoID:        eventoModel.ID,
			Nombre:          tipoTicket.Label,
			FechaIni:        fechaIni,
			FechaFin:        fechaFin,
			Estado:          1,
			UsuarioCreacion: &usuarioCreacion,
			FechaCreacion:   now,
		}
		if err := tx.Create(tipoTicketModel).Error; err != nil {
			tx.Rollback()
			e.logger.Errorf("Failed to create tipo de ticket: %v", err)
			return nil, &errors.BadRequestError.EventoNotCreated
		}
		tiposTicketMap[tipoTicket.ID] = tipoTicketModel.ID
	}

	// Create tarifas (prices) based on the pricing matrix
	for sectorID, perfilesPrecios := range eventoReq.Precios {
		sectorDBID, ok := sectoresMap[sectorID]
		if !ok {
			tx.Rollback()
			e.logger.Errorf("Sector ID %s not found in sectoresMap", sectorID)
			return nil, &errors.BadRequestError.EventoNotCreated
		}

		for perfilID, tiposTicketPrecios := range perfilesPrecios {
			perfilDBID, ok := perfilesMap[perfilID]
			if !ok {
				tx.Rollback()
				e.logger.Errorf("Perfil ID %s not found in perfilesMap", perfilID)
				return nil, &errors.BadRequestError.EventoNotCreated
			}

			for tipoTicketID, precio := range tiposTicketPrecios {
				tipoTicketDBID, ok := tiposTicketMap[tipoTicketID]
				if !ok {
					tx.Rollback()
					e.logger.Errorf("Tipo ticket ID %s not found in tiposTicketMap", tipoTicketID)
					return nil, &errors.BadRequestError.EventoNotCreated
				}

				tarifaModel := &model.Tarifa{
					SectorID:          sectorDBID,
					TipoDeTicketID:    tipoTicketDBID,
					PerfilDePersonaID: &perfilDBID,
					Precio:            precio,
					Estado:            1,
					UsuarioCreacion:   &usuarioCreacion,
					FechaCreacion:     now,
				}
				if err := tx.Create(tarifaModel).Error; err != nil {
					tx.Rollback()
					e.logger.Errorf("Failed to create tarifa: %v", err)
					return nil, &errors.BadRequestError.EventoNotCreated
				}
			}
		}
	}

	// Create event dates
	eventDatesResponse := []schemas.EventDateResponse{}
	for _, eventDate := range eventoReq.EventDates {
		// Parse the date
		fecha, err := time.Parse("2006-01-02", eventDate.Fecha)
		if err != nil {
			tx.Rollback()
			e.logger.Errorf("Failed to parse fecha: %v", err)
			return nil, &errors.BadRequestError.EventoNotCreated
		}

		// Check if fecha already exists or create new one
		var fechaModel model.Fecha
		result := tx.Where("fecha_evento = ?", fecha).First(&fechaModel)
		if result.Error == gorm.ErrRecordNotFound {
			// Create new fecha
			fechaModel = model.Fecha{
				FechaEvento: fecha,
			}
			if err := tx.Create(&fechaModel).Error; err != nil {
				tx.Rollback()
				e.logger.Errorf("Failed to create fecha: %v", err)
				return nil, &errors.BadRequestError.EventoNotCreated
			}
		}

		// Parse hora inicio
		horaInicio, err := time.Parse("15:04", eventDate.HoraInicio)
		if err != nil {
			tx.Rollback()
			e.logger.Errorf("Failed to parse hora inicio: %v", err)
			return nil, &errors.BadRequestError.EventoNotCreated
		}
		// Combine date and time
		horaInicioFull := time.Date(fecha.Year(), fecha.Month(), fecha.Day(),
			horaInicio.Hour(), horaInicio.Minute(), 0, 0, time.UTC)

		// Create evento_fecha
		eventoFechaModel := &model.EventoFecha{
			EventoID:        eventoModel.ID,
			FechaID:         fechaModel.ID,
			HoraInicio:      horaInicioFull,
			Estado:          1,
			UsuarioCreacion: &usuarioCreacion,
			FechaCreacion:   now,
		}
		if err := tx.Create(eventoFechaModel).Error; err != nil {
			tx.Rollback()
			e.logger.Errorf("Failed to create evento_fecha: %v", err)
			return nil, &errors.BadRequestError.EventoNotCreated
		}

		eventDatesResponse = append(eventDatesResponse, schemas.EventDateResponse{
			IdFechaEvento: eventoFechaModel.ID,
			IdFecha:       fechaModel.ID,
			Fecha:         eventDate.Fecha,
			HoraInicio:    eventDate.HoraInicio,
			HoraFin:       eventDate.HoraFin,
		})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		e.logger.Errorf("Failed to commit transaction: %v", err)
		return nil, &errors.BadRequestError.EventoNotCreated
	}

	// Convert request types to response types
	perfilesResponse := make([]schemas.PerfilResponse, len(eventoReq.Perfiles))
	for i, p := range eventoReq.Perfiles {
		perfilesResponse[i] = schemas.PerfilResponse{
			ID:    p.ID,
			Label: p.Label,
		}
	}

	sectoresResponse := make([]schemas.SectorResponse, len(eventoReq.Sectores))
	for i, s := range eventoReq.Sectores {
		sectoresResponse[i] = schemas.SectorResponse{
			ID:        s.ID,
			Nombre:    s.Nombre,
			Capacidad: s.Capacidad,
		}
	}

	tiposTicketResponse := make([]schemas.TipoTicketResponse, len(eventoReq.TiposTicket))
	for i, t := range eventoReq.TiposTicket {
		tiposTicketResponse[i] = schemas.TipoTicketResponse{
			ID:    t.ID,
			Label: t.Label,
		}
	}

	// Build the response
	response := &schemas.EventoResponse{
		IdEvento:          eventoModel.ID,
		IdOrganizador:     eventoModel.OrganizadorID,
		IdCategoria:       eventoModel.CategoriaID,
		Titulo:            eventoModel.Titulo,
		Descripcion:       eventoModel.Descripcion,
		Lugar:             eventoModel.Lugar,
		Estado:            convert.MapEstadoToString(eventoModel.EventoEstado),
		Likes:             eventoModel.CantMeGusta,
		NoInteres:         eventoModel.CantNoInteresa,
		CantVendidasTotal: eventoModel.CantVendidoTotal,
		TotalRecaudado:    eventoModel.TotalRecaudado,
		ImagenPortada:     eventoModel.ImagenPortada,
		ImagenLugar:       eventoModel.ImagenEscenario,
		VideoUrl:          eventoModel.VideoPresentacion,
		EventDates:        eventDatesResponse,
		Perfiles:          perfilesResponse,
		Sectores:          sectoresResponse,
		TiposTicket:       tiposTicketResponse,
		Precios:           eventoReq.Precios,
		Metadata: schemas.MetadataResponse{
			CreadoPor:           eventoReq.Metadata.CreadoPor,
			FechaCreacion:       eventoModel.FechaCreacion.Format(time.RFC3339),
			UltimaActualizacion: eventoModel.FechaCreacion.Format(time.RFC3339),
			Version:             eventoReq.Metadata.Version,
		},
	}

	return response, nil
}

// EditarEventoFull reemplaza completamente un evento (solo BORRADOR y sin ventas).
// Borra dependencias y las recrea con el mismo formato de creación.
func (e *Evento) EditarEventoFull(eventoID int64, req *schemas.EditarEventoFullRequest) (*schemas.EventoResponse, *errors.Error) {
	if eventoID <= 0 {
		return nil, &errors.BadRequestError.InvalidIDParam
	}

	tx := e.DaoPostgresql.Evento.PostgresqlDB.Begin()
	if tx.Error != nil {
		e.logger.Errorf("EditarEventoFull begin tx: %v", tx.Error)
		return nil, &errors.InternalServerError.Default
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()

	var ev model.Evento
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&ev, "evento_id = ?", eventoID).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}
		e.logger.Errorf("EditarEventoFull fetch evento id=%d: %v", eventoID, err)
		return nil, &errors.InternalServerError.Default
	}

	// Validaciones: solo borrador y sin ventas/tickets.
	if ev.EventoEstado != convert.MapEstadoToInt16("BORRADOR") {
		tx.Rollback()
		return nil, &errors.UnprocessableEntityError.InvalidRequestBody
	}
	if ev.CantVendidoTotal > 0 {
		tx.Rollback()
		return nil, &errors.UnprocessableEntityError.InvalidRequestBody
	}
	var ticketsCount int64
	if err := tx.Model(&model.Ticket{}).
		Joins("JOIN evento_fecha ef ON ef.evento_fecha_id = ticket.evento_fecha_id").
		Where("ef.evento_id = ?", eventoID).
		Count(&ticketsCount).Error; err != nil {
		tx.Rollback()
		e.logger.Errorf("EditarEventoFull count tickets evento=%d: %v", eventoID, err)
		return nil, &errors.InternalServerError.Default
	}
	if ticketsCount > 0 {
		tx.Rollback()
		return nil, &errors.UnprocessableEntityError.InvalidRequestBody
	}

	// Actualizar cabecera de evento
	ev.OrganizadorID = req.IdOrganizador
	ev.CategoriaID = req.IdCategoria
	ev.Titulo = req.Titulo
	ev.Descripcion = req.Descripcion
	ev.Lugar = req.Lugar
	ev.EventoEstado = convert.MapEstadoToInt16(req.Estado)
	ev.CantMeGusta = req.Likes
	ev.CantNoInteresa = req.NoInteres
	ev.CantVendidoTotal = req.CantVendidasTotal
	ev.TotalRecaudado = req.TotalRecaudado
	ev.ImagenPortada = req.ImagenPortada
	ev.ImagenEscenario = req.ImagenLugar
	ev.VideoPresentacion = req.VideoUrl
	ev.UsuarioModificacion = &req.UsuarioModificacion
	ev.FechaModificacion = &now

	if err := tx.Save(&ev).Error; err != nil {
		tx.Rollback()
		e.logger.Errorf("EditarEventoFull update evento id=%d: %v", eventoID, err)
		return nil, &errors.InternalServerError.Default
	}

	// Eliminar dependencias
	if err := tx.Where("evento_fecha_id IN (?)",
		tx.Table("evento_fecha").Select("evento_fecha_id").Where("evento_id = ?", eventoID)).
		Delete(&model.Ticket{}).Error; err != nil {
		tx.Rollback()
		e.logger.Errorf("EditarEventoFull delete tickets evento=%d: %v", eventoID, err)
		return nil, &errors.InternalServerError.Default
	}
	if err := tx.Where("sector_id IN (?)",
		tx.Table("sector").Select("sector_id").Where("evento_id = ?", eventoID)).
		Delete(&model.Tarifa{}).Error; err != nil {
		tx.Rollback()
		e.logger.Errorf("EditarEventoFull delete tarifas-sector evento=%d: %v", eventoID, err)
		return nil, &errors.InternalServerError.Default
	}
	if err := tx.Where("tipo_de_ticket_id IN (?)",
		tx.Table("tipo_de_ticket").Select("tipo_de_ticket_id").Where("evento_id = ?", eventoID)).
		Delete(&model.Tarifa{}).Error; err != nil {
		tx.Rollback()
		e.logger.Errorf("EditarEventoFull delete tarifas-tipo evento=%d: %v", eventoID, err)
		return nil, &errors.InternalServerError.Default
	}
	if err := tx.Where("perfil_de_persona_id IN (?)",
		tx.Table("perfil_de_persona").Select("perfil_de_persona_id").Where("evento_id = ?", eventoID)).
		Delete(&model.Tarifa{}).Error; err != nil {
		tx.Rollback()
		e.logger.Errorf("EditarEventoFull delete tarifas-perfil evento=%d: %v", eventoID, err)
		return nil, &errors.InternalServerError.Default
	}

	if err := tx.Where("evento_id = ?", eventoID).Delete(&model.EventoFecha{}).Error; err != nil {
		tx.Rollback()
		e.logger.Errorf("EditarEventoFull delete evento_fecha evento=%d: %v", eventoID, err)
		return nil, &errors.InternalServerError.Default
	}
	if err := tx.Where("evento_id = ?", eventoID).Delete(&model.Sector{}).Error; err != nil {
		tx.Rollback()
		e.logger.Errorf("EditarEventoFull delete sector evento=%d: %v", eventoID, err)
		return nil, &errors.InternalServerError.Default
	}
	if err := tx.Where("evento_id = ?", eventoID).Delete(&model.PerfilDePersona{}).Error; err != nil {
		tx.Rollback()
		e.logger.Errorf("EditarEventoFull delete perfil evento=%d: %v", eventoID, err)
		return nil, &errors.InternalServerError.Default
	}
	if err := tx.Where("evento_id = ?", eventoID).Delete(&model.TipoDeTicket{}).Error; err != nil {
		tx.Rollback()
		e.logger.Errorf("EditarEventoFull delete tipo_ticket evento=%d: %v", eventoID, err)
		return nil, &errors.InternalServerError.Default
	}

	usuario := req.UsuarioModificacion

	// Perfiles
	perfilesMap := make(map[string]int64)
	for _, perfil := range req.Perfiles {
		perfilModel := &model.PerfilDePersona{
			EventoID:        ev.ID,
			Nombre:          perfil.Label,
			Estado:          1,
			UsuarioCreacion: &usuario,
			FechaCreacion:   now,
		}
		if err := tx.Create(perfilModel).Error; err != nil {
			tx.Rollback()
			e.logger.Errorf("EditarEventoFull create perfil: %v", err)
			return nil, &errors.BadRequestError.EventoNotCreated
		}
		perfilesMap[perfil.ID] = perfilModel.ID
	}

	// Sectores
	sectoresMap := make(map[string]int64)
	for _, sector := range req.Sectores {
		sectorModel := &model.Sector{
			EventoID:        ev.ID,
			SectorTipo:      sector.Nombre,
			TotalEntradas:   sector.Capacidad,
			CantVendidas:    0,
			Estado:          1,
			UsuarioCreacion: &usuario,
			FechaCreacion:   now,
		}
		if err := tx.Create(sectorModel).Error; err != nil {
			tx.Rollback()
			e.logger.Errorf("EditarEventoFull create sector: %v", err)
			return nil, &errors.BadRequestError.EventoNotCreated
		}
		sectoresMap[sector.ID] = sectorModel.ID
	}

	// Tipos de ticket
	tiposTicketMap := make(map[string]int64)
	for _, tipoTicket := range req.TiposTicket {
		fechaIni := now
		fechaFin := now.AddDate(1, 0, 0)
		tipoTicketModel := &model.TipoDeTicket{
			EventoID:        ev.ID,
			Nombre:          tipoTicket.Label,
			FechaIni:        fechaIni,
			FechaFin:        fechaFin,
			Estado:          1,
			UsuarioCreacion: &usuario,
			FechaCreacion:   now,
		}
		if err := tx.Create(tipoTicketModel).Error; err != nil {
			tx.Rollback()
			e.logger.Errorf("EditarEventoFull create tipo_ticket: %v", err)
			return nil, &errors.BadRequestError.EventoNotCreated
		}
		tiposTicketMap[tipoTicket.ID] = tipoTicketModel.ID
	}

	// Tarifas
	for sectorID, perfilesPrecios := range req.Precios {
		sectorDBID, ok := sectoresMap[sectorID]
		if !ok {
			tx.Rollback()
			e.logger.Errorf("EditarEventoFull sector not found %s", sectorID)
			return nil, &errors.BadRequestError.EventoNotCreated
		}

		for perfilID, tiposTicketPrecios := range perfilesPrecios {
			perfilDBID, ok := perfilesMap[perfilID]
			if !ok {
				tx.Rollback()
				e.logger.Errorf("EditarEventoFull perfil not found %s", perfilID)
				return nil, &errors.BadRequestError.EventoNotCreated
			}

			for tipoTicketID, precio := range tiposTicketPrecios {
				tipoTicketDBID, ok := tiposTicketMap[tipoTicketID]
				if !ok {
					tx.Rollback()
					e.logger.Errorf("EditarEventoFull tipo_ticket not found %s", tipoTicketID)
					return nil, &errors.BadRequestError.EventoNotCreated
				}

				tarifaModel := &model.Tarifa{
					SectorID:          sectorDBID,
					TipoDeTicketID:    tipoTicketDBID,
					PerfilDePersonaID: &perfilDBID,
					Precio:            precio,
					Estado:            1,
					UsuarioCreacion:   &usuario,
					FechaCreacion:     now,
				}
				if err := tx.Create(tarifaModel).Error; err != nil {
					tx.Rollback()
					e.logger.Errorf("EditarEventoFull create tarifa: %v", err)
					return nil, &errors.BadRequestError.EventoNotCreated
				}
			}
		}
	}

	// Fechas del evento
	eventDatesResponse := []schemas.EventDateResponse{}
	for _, eventDate := range req.EventDates {
		fecha, err := time.Parse("2006-01-02", eventDate.Fecha)
		if err != nil {
			tx.Rollback()
			e.logger.Errorf("EditarEventoFull parse fecha: %v", err)
			return nil, &errors.BadRequestError.EventoNotCreated
		}

		var fechaModel model.Fecha
		result := tx.Where("fecha_evento = ?", fecha).First(&fechaModel)
		if result.Error == gorm.ErrRecordNotFound {
			fechaModel = model.Fecha{FechaEvento: fecha}
			if err := tx.Create(&fechaModel).Error; err != nil {
				tx.Rollback()
				e.logger.Errorf("EditarEventoFull create fecha: %v", err)
				return nil, &errors.BadRequestError.EventoNotCreated
			}
		} else if result.Error != nil {
			tx.Rollback()
			e.logger.Errorf("EditarEventoFull find fecha: %v", result.Error)
			return nil, &errors.BadRequestError.EventoNotCreated
		}

		horaInicio, err := time.Parse("15:04", eventDate.HoraInicio)
		if err != nil {
			tx.Rollback()
			e.logger.Errorf("EditarEventoFull parse hora inicio: %v", err)
			return nil, &errors.BadRequestError.EventoNotCreated
		}
		horaInicioFull := time.Date(fecha.Year(), fecha.Month(), fecha.Day(),
			horaInicio.Hour(), horaInicio.Minute(), 0, 0, time.UTC)

		eventoFechaModel := &model.EventoFecha{
			EventoID:        ev.ID,
			FechaID:         fechaModel.ID,
			HoraInicio:      horaInicioFull,
			Estado:          1,
			UsuarioCreacion: &usuario,
			FechaCreacion:   now,
		}
		if err := tx.Create(eventoFechaModel).Error; err != nil {
			tx.Rollback()
			e.logger.Errorf("EditarEventoFull create evento_fecha: %v", err)
			return nil, &errors.BadRequestError.EventoNotCreated
		}

		eventDatesResponse = append(eventDatesResponse, schemas.EventDateResponse{
			IdFechaEvento: eventoFechaModel.ID,
			IdFecha:       fechaModel.ID,
			Fecha:         eventDate.Fecha,
			HoraInicio:    eventDate.HoraInicio,
			HoraFin:       eventDate.HoraFin,
		})
	}

	if err := tx.Commit().Error; err != nil {
		e.logger.Errorf("EditarEventoFull commit evento=%d: %v", eventoID, err)
		return nil, &errors.BadRequestError.EventoNotCreated
	}

	// Construir response
	perfilesResponse := make([]schemas.PerfilResponse, len(req.Perfiles))
	for i, p := range req.Perfiles {
		perfilesResponse[i] = schemas.PerfilResponse{
			ID:    p.ID,
			Label: p.Label,
		}
	}

	sectoresResponse := make([]schemas.SectorResponse, len(req.Sectores))
	for i, s := range req.Sectores {
		sectoresResponse[i] = schemas.SectorResponse{
			ID:        s.ID,
			Nombre:    s.Nombre,
			Capacidad: s.Capacidad,
		}
	}

	tiposTicketResponse := make([]schemas.TipoTicketResponse, len(req.TiposTicket))
	for i, t := range req.TiposTicket {
		tiposTicketResponse[i] = schemas.TipoTicketResponse{
			ID:    t.ID,
			Label: t.Label,
		}
	}

	creadoPor := req.Metadata.CreadoPor
	if creadoPor == "" && ev.UsuarioCreacion != nil {
		creadoPor = fmt.Sprintf("%d", *ev.UsuarioCreacion)
	}

	response := &schemas.EventoResponse{
		IdEvento:          ev.ID,
		IdOrganizador:     ev.OrganizadorID,
		IdCategoria:       ev.CategoriaID,
		Titulo:            ev.Titulo,
		Descripcion:       ev.Descripcion,
		Lugar:             ev.Lugar,
		Estado:            convert.MapEstadoToString(ev.EventoEstado),
		Likes:             ev.CantMeGusta,
		NoInteres:         ev.CantNoInteresa,
		CantVendidasTotal: ev.CantVendidoTotal,
		TotalRecaudado:    ev.TotalRecaudado,
		ImagenPortada:     ev.ImagenPortada,
		ImagenLugar:       ev.ImagenEscenario,
		VideoUrl:          ev.VideoPresentacion,
		EventDates:        eventDatesResponse,
		Perfiles:          perfilesResponse,
		Sectores:          sectoresResponse,
		TiposTicket:       tiposTicketResponse,
		Precios:           req.Precios,
		Metadata: schemas.MetadataResponse{
			CreadoPor:           creadoPor,
			FechaCreacion:       ev.FechaCreacion.Format(time.RFC3339),
			UltimaActualizacion: now.Format(time.RFC3339),
			Version:             req.Metadata.Version,
		},
	}

	return response, nil
}

// FetchPostgresqlEventos retrieves the upcoming events without filters
func (e *Evento) FetchPostgresqlEventos() (*schemas.EventosPaginados, *errors.Error) {
	eventos, err := e.DaoPostgresql.Evento.ObtenerEventosDisponiblesSinFiltros()
	if err != nil {
		e.logger.Errorf("Failed to fetch eventos: %v", err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	mapEventDates(eventos)

	total := int64(len(eventos))
	totalPaginas := 0
	if total > 0 {
		totalPaginas = 1
	}

	return &schemas.EventosPaginados{
		Eventos:      eventos,
		Total:        total,
		PaginaActual: 1,
		TotalPaginas: totalPaginas,
	}, nil
}

// FetchPostgresqlEventosFeed obtiene los eventos para el feed de recomendaciones
func (e *Evento) FetchPostgresqlEventosFeed(usuarioId *int64) (*schemas.EventosPaginados, *errors.Error) {
	eventos, err := e.DaoPostgresql.Evento.ObtenerEventosParaElFeed(usuarioId)
	if err != nil {
		e.logger.Errorf("Failed to fetch eventos: %v", err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	mapEventDates(eventos)

	total := int64(len(eventos))
	totalPaginas := 0
	if total > 0 {
		totalPaginas = 1
	}

	return &schemas.EventosPaginados{
		Eventos:      eventos,
		Total:        total,
		PaginaActual: 1,
		TotalPaginas: totalPaginas,
	}, nil
}
func (e *Evento) FetchPostgresqlEventosConInteraccionesFeed(usuarioId *int64) (*schemas.EventosPaginados, *errors.Error) {
	eventos, err := e.DaoPostgresql.Evento.CargarEventosNuevamenteParaElFeed(usuarioId)
	if err != nil {
		e.logger.Errorf("Failed to fetch eventos: %v", err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	mapEventDates(eventos)

	total := int64(len(eventos))
	totalPaginas := 0
	if total > 0 {
		totalPaginas = 1
	}

	return &schemas.EventosPaginados{
		Eventos:      eventos,
		Total:        total,
		PaginaActual: 1,
		TotalPaginas: totalPaginas,
	}, nil
}

// FetchPostgresqlEventos retrieves the upcoming events with filters
func (e *Evento) FetchPostgresqlEventosWithFilters(
	categoriaID *int64,
	organizadorID *int64,
	titulo *string,
	descripcion *string,
	lugar *string,
	fecha *time.Time,
	horaInicio *time.Time,
	estado *int16,
	soloFuturos bool) (*schemas.EventosPaginados, *errors.Error) {
	eventos, err := e.DaoPostgresql.Evento.ObtenerEventosDisponiblesConFiltros(
		categoriaID,
		organizadorID,
		titulo,
		descripcion,
		lugar,
		fecha,
		horaInicio,
		estado,
		soloFuturos)

	if err != nil {
		e.logger.Errorf("Failed to fetch eventos: %v", err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	mapEventDates(eventos)

	//Revisar esta lógica
	total := int64(len(eventos))
	totalPaginas := 0
	if total > 0 {
		totalPaginas = 1
	}

	return &schemas.EventosPaginados{
		Eventos:      eventos,
		Total:        total,
		PaginaActual: 1,
		TotalPaginas: totalPaginas,
	}, nil
}

// mapEventDates rellena el slice EventDates con fechas formateadas para la respuesta JSON.
func mapEventDates(eventos []*model.Evento) {
	for _, ev := range eventos {
		if ev == nil {
			continue
		}
		ev.EventDates = make([]model.EventDateView, 0, len(ev.Fechas))
		for _, ef := range ev.Fechas {
			fechaStr := ""
			if ef.Fecha != nil {
				fechaStr = ef.Fecha.FechaEvento.Format("2006-01-02")
			}
			horaInicioStr := ef.HoraInicio.Format("15:04")

			ev.EventDates = append(ev.EventDates, model.EventDateView{
				IdFechaEvento: ef.ID,
				IdFecha:       ef.FechaID,
				Fecha:         fechaStr,
				HoraInicio:    horaInicioStr,
				HoraFin:       "",
			})
		}
	}
}

// GetPostgresqlEventoById gets an event by ID
func (e *Evento) GetPostgresqlEventoById(eventoID int64) (*schemas.EventoResponse, *errors.Error) {
	var eventoModel model.Evento

	// Use preload to fetch all related entities
	result := e.DaoPostgresql.Evento.PostgresqlDB.
		Preload("Perfiles").
		Preload("Sectores").
		Preload("TiposTicket").
		Preload("Fechas.Fecha").
		First(&eventoModel, eventoID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}
		e.logger.Errorf("Failed to get evento: %v", result.Error)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	// Build event dates response
	eventDates := []schemas.EventDateResponse{}
	for _, ef := range eventoModel.Fechas {
		if ef.Fecha != nil {
			eventDates = append(eventDates, schemas.EventDateResponse{
				IdFechaEvento: ef.ID,
				IdFecha:       ef.FechaID,
				Fecha:         ef.Fecha.FechaEvento.Format("2006-01-02"),
				HoraInicio:    ef.HoraInicio.Format("15:04"),
				HoraFin:       "", // You may need to add HoraFin to the model
			})
		}
	}

	// Build perfiles response
	perfiles := []schemas.PerfilResponse{}
	for _, p := range eventoModel.Perfiles {
		perfiles = append(perfiles, schemas.PerfilResponse{
			ID:    fmt.Sprintf("%d", p.ID),
			Label: p.Nombre,
		})
	}

	// Build sectores response
	sectores := []schemas.SectorResponse{}
	for _, s := range eventoModel.Sectores {
		sectores = append(sectores, schemas.SectorResponse{
			ID:        fmt.Sprintf("%d", s.ID),
			Nombre:    s.SectorTipo,
			Capacidad: s.TotalEntradas,
		})
	}

	// Build tipos ticket response
	tiposTicket := []schemas.TipoTicketResponse{}
	for _, t := range eventoModel.TiposTicket {
		tiposTicket = append(tiposTicket, schemas.TipoTicketResponse{
			ID:    fmt.Sprintf("%d", t.ID),
			Label: t.Nombre,
		})
	}

	// Fetch and build precios desde las tarifas activas del evento
	tarifas, tarifaErr := e.DaoPostgresql.Tarifa.ListarTarifasPorEvento(eventoID)
	if tarifaErr != nil {
		e.logger.Errorf("Failed to list tarifas for evento %d: %v", eventoID, tarifaErr)
		return nil, &errors.InternalServerError.Default
	}

	precios := make(schemas.PreciosSector)
	for _, tarifa := range tarifas {
		sectorKey := fmt.Sprintf("%d", tarifa.SectorID)

		perfilKey := "0"
		if tarifa.PerfilDePersonaID != nil {
			perfilKey = fmt.Sprintf("%d", *tarifa.PerfilDePersonaID)
		}

		tipoTicketKey := fmt.Sprintf("%d", tarifa.TipoDeTicketID)

		perfilMap, ok := precios[sectorKey]
		if !ok {
			perfilMap = make(schemas.PreciosPerfil)
			precios[sectorKey] = perfilMap
		}

		precioDetalle, ok := perfilMap[perfilKey]
		if !ok {
			precioDetalle = make(schemas.PrecioDetalle)
			perfilMap[perfilKey] = precioDetalle
		}

		precioDetalle[tipoTicketKey] = tarifa.Precio
	}

	response := &schemas.EventoResponse{
		IdEvento:          eventoModel.ID,
		IdOrganizador:     eventoModel.OrganizadorID,
		IdCategoria:       eventoModel.CategoriaID,
		Titulo:            eventoModel.Titulo,
		Descripcion:       eventoModel.Descripcion,
		Lugar:             eventoModel.Lugar,
		Estado:            convert.MapEstadoToString(eventoModel.EventoEstado),
		Likes:             eventoModel.CantMeGusta,
		NoInteres:         eventoModel.CantNoInteresa,
		CantVendidasTotal: eventoModel.CantVendidoTotal,
		TotalRecaudado:    eventoModel.TotalRecaudado,
		ImagenPortada:     eventoModel.ImagenPortada,
		ImagenLugar:       eventoModel.ImagenEscenario,
		VideoUrl:          eventoModel.VideoPresentacion,
		EventDates:        eventDates,
		Perfiles:          perfiles,
		Sectores:          sectores,
		TiposTicket:       tiposTicket,
		Precios:           precios,
		Metadata: schemas.MetadataResponse{
			CreadoPor:           fmt.Sprintf("%d", *eventoModel.UsuarioCreacion),
			FechaCreacion:       eventoModel.FechaCreacion.Format(time.RFC3339),
			UltimaActualizacion: eventoModel.FechaCreacion.Format(time.RFC3339),
			Version:             1,
		},
	}

	return response, nil
}

func (e *Evento) GetPostgresqlReporteEvento(
	organizadorID int64,
	eventoID *int64,
	fechaDesde *time.Time,
	fechaHasta *time.Time,
) ([]*schemas.EventoReporte, *errors.Error) {
	var eventos []*model.Evento

	if eventoID != nil {
		evento, err := e.DaoPostgresql.Evento.ObtenerEventoPorId(*eventoID)
		if err != nil {
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}

		eventos = append(eventos, evento)
	} else {
		var err error
		eventos, err = e.DaoPostgresql.Evento.ObtenerEventosDelOrganizador(organizadorID)
		if err != nil {
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}
	}

	eventoReporte := []*schemas.EventoReporte{}
	for _, ev := range eventos {
		capacidadEvento, _ := e.DaoPostgresql.Sector.ObtenerCapacidadPorEvento(ev.ID)
		ingresoEvento, cargos, ticketVendido := e.DaoPostgresql.OrdenDeCompra.ObtenerIngresoCargoPorFecha(ev.ID, fechaDesde, fechaHasta)

		eventoReporte = append(eventoReporte, &schemas.EventoReporte{
			IdEvento:         ev.ID,
			Titulo:           ev.Titulo,
			Lugar:            ev.Lugar,
			Capacidad:        capacidadEvento,                 //calcular capacidad con sector
			IngresoTotal:     ingresoEvento,                   //calcular con orden de compra
			TicketsVendidos:  ticketVendido,                   //calcular con orden de compra
			CargosPorServico: cargos,                          //calcular con orden de compra
			Comisiones:       (ingresoEvento - cargos) * 0.05, //(Ingreso total - cargo)*5%
			VentasPorTipo:    []schemas.TipoTicketReporte{},
			Fechas:           []schemas.EventDateReporte{},
		})
	}

	return eventoReporte, nil
}

// reporte mamadisimo del organizador
func (e *Evento) GetPostgresqlReporteEventosOrganizador(
	organizadorID int64,
	fechaDesde *time.Time,
	fechaHasta *time.Time,
) ([]schemas.EventoOrganizadorReporte, *errors.Error) {
	eventos, err := e.DaoPostgresql.Evento.ObtenerEventosDelOrganizador(organizadorID)
	if err != nil {
		e.logger.Errorf("Failed to fetch eventos for organizer %d: %v", organizadorID, err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	reportes := make([]schemas.EventoOrganizadorReporte, 0, len(eventos))

	for _, ev := range eventos {
		capacidadEvento, capErr := e.DaoPostgresql.Sector.ObtenerCapacidadPorEvento(ev.ID)
		if capErr != nil {
			e.logger.Warnf("Capacidad no disponible para evento %d: %v", ev.ID, capErr)
			capacidadEvento = 0
		}

		ingresoTotal, cargosServicio, ticketsVendidos := e.DaoPostgresql.OrdenDeCompra.ObtenerIngresoCargoPorFecha(ev.ID, fechaDesde, fechaHasta)
		ventasPorSectorDTO, ventasErr := e.DaoPostgresql.OrdenDeCompra.ObtenerVentasPorSector(ev.ID, fechaDesde, fechaHasta)
		if ventasErr != nil {
			ventasPorSectorDTO = []daoPostgresql.VentaPorSectorDTO{}
		}

		ventasPorSector := make([]schemas.VentaPorSectorOrganizador, 0, len(ventasPorSectorDTO))
		for _, v := range ventasPorSectorDTO {
			ventasPorSector = append(ventasPorSector, schemas.VentaPorSectorOrganizador{
				Sector:   v.Sector,
				Vendidos: v.TicketsVendidos,
				Ingresos: v.Ingresos,
			})
		}

		fechas := make([]schemas.EventoFechaOrganizadorReporte, 0, len(ev.Fechas))
		for _, f := range ev.Fechas {
			fechaStr := ""
			if f.Fecha != nil {
				fechaStr = f.Fecha.FechaEvento.Format("2006-01-02")
			}

			horaInicio := ""
			if !f.HoraInicio.IsZero() {
				horaInicio = f.HoraInicio.Format("15:04")
			}

			fechas = append(fechas, schemas.EventoFechaOrganizadorReporte{
				IdFechaEvento: f.ID,
				Fecha:         fechaStr,
				HoraInicio:    horaInicio,
				HoraFin:       "",
			})
		}

		estado := deriveEstadoEventoOrganizador(ev.EventoEstado, capacidadEvento, ticketsVendidos)
		gananciaNeta := ingresoTotal - cargosServicio
		comisiones := gananciaNeta * 0.05
		if gananciaNeta < 0 {
			gananciaNeta = 0 // por si acaso, evitar negativos raros
		}
		reportes = append(reportes, schemas.EventoOrganizadorReporte{
			IdEvento:        ev.ID,
			Nombre:          ev.Titulo,
			Ubicacion:       ev.Lugar,
			Capacidad:       capacidadEvento,
			Estado:          estado,
			IngresosTotales: ingresoTotal,
			GananciaNeta:    gananciaNeta,
			TicketsVendidos: ticketsVendidos,
			VentasPorSector: ventasPorSector,
			Fechas:          fechas,
			CargosServicio:  cargosServicio,
			Comisiones:      comisiones,
		})
	}

	return reportes, nil
}

func deriveEstadoEventoOrganizador(eventoEstado int16, capacidad int64, ticketsVendidos int64) string {
	if eventoEstado == convert.MapEstadoToInt16("CANCELADO") {
		return "CANCELADO"
	}

	if capacidad > 0 && ticketsVendidos >= capacidad {
		return "AGOTADO"
	}

	if eventoEstado == convert.MapEstadoToInt16("PUBLICADO") {
		return "EN_VENTA"
	}

	return "BORRADOR"
}

func (a *Evento) GenerarReporteAdministrativo(req schemas.AdminReportRequest) (*schemas.AdminReportResponse, *errors.Error) {

	// 1. Validaciones y Defaults
	limit := 100
	if req.Limit > 0 {
		limit = req.Limit
	}

	// 2. Conversión de Fechas (ISO String -> Time)
	var fechaInicio, fechaFin *time.Time
	if req.FechaInicio != nil && *req.FechaInicio != "" {
		t, err := time.Parse(time.RFC3339, *req.FechaInicio)
		if err != nil {
			return nil, &errors.UnprocessableEntityError.InvalidDateFormat
		}
		fechaInicio = &t
	}
	if req.FechaFin != nil && *req.FechaFin != "" {
		t, err := time.Parse(time.RFC3339, *req.FechaFin)
		if err != nil {
			return nil, &errors.UnprocessableEntityError.InvalidDateFormat
		}
		fechaFin = &t
	}

	// 3. Conversión de Estado (String -> Int16 para DB)
	// Asumiendo: 0=BORRADOR, 1=PUBLICADO, 2=CANCELADO (Ajusta según tu modelo real)
	var estadoInt *int16
	if req.Estado != "" {
		val := convert.MapEstadoToInt16(req.Estado) // Tu función utilitaria existente
		estadoInt = &val
	}

	// 4. Llamada al DAO
	reporte, err := a.DaoPostgresql.Evento.GenerarReporteAdmin(
		fechaInicio,
		fechaFin,
		req.IdCategoria,
		req.IdOrganizador,
		estadoInt,
		limit,
	)

	if err != nil {
		a.logger.Errorf("GenerarReporteAdmin Error: %v", err)
		return nil, &errors.InternalServerError.Default
	}

	// 5. Manejo de "Sin Datos" (204/404)
	if reporte == nil {
		// Esto retornará JSON { "code": "NO_DATA", "message": "..." } con Status 404
		return nil, &errors.ObjectNotFoundError.ReportNoDataFound
	}

	// 6. Mapeo Final de Estados (Int -> String) para la lista de eventos
	// El DAO devolvió el int, ahora lo pasamos a string para el JSON
	for i := range reporte.Events {
		// Aquí asumimos que en DAO escaneaste evento_estado en un campo temporal o usas el helper
		// Si el scan falló en mapear string directo, hazlo aquí:
		// reporte.Events[i].Estado = convert.MapEstadoToString(reporte.Events[i].EstadoInt)

		// Como ejemplo simple, si ya viene string o lo mapeas:
		if reporte.Events[i].Estado == "" {
			reporte.Events[i].Estado = "PUBLICADO" // Fallback o lógica real
		}
	}

	return reporte, nil
}

func (e *Evento) GetPostgresqlEventoDetalle(eventoId int64) (*schemas.EventoDetalleDTO, *errors.Error) {
	eventoDetalle, err := e.DaoPostgresql.Evento.ObtenerEventoDetalle(eventoId)
	if err != nil {
		e.logger.Errorf("Failed to get evento detallado: %v", err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	return eventoDetalle, nil
}

// EditarEvento aplica cambios a:
// - Evento (ubicación, estados)
// - Fechas (fecha calendario, hora inicio, reasignar fecha_id)
// - Sectores
// - Perfiles de persona
// - Tipos de ticket
// y devuelve el detalle actualizado del evento.
func (e *Evento) EditarEvento(
	req *schemas.EditarEventoRequest,
) (*schemas.EventoDetalleDTO, *errors.Error) {

	if req.IdEvento <= 0 {
		return nil, &errors.BadRequestError.InvalidIDParam
	}

	now := time.Now()
	userID := req.UsuarioModificacion

	updates := map[string]any{}
	if req.NuevaDescripcion != nil && *req.NuevaDescripcion != "" {
		updates["descripcion"] = *req.NuevaDescripcion
	}
	if req.NuevaImagenPortada != nil && *req.NuevaImagenPortada != "" {
		updates["imagen_portada"] = *req.NuevaImagenPortada
	}
	if req.NuevaImagenPresentacion != nil && *req.NuevaImagenPresentacion != "" {
		updates["imagen_escenario"] = *req.NuevaImagenPresentacion
	}
	if req.NuevoVideo != nil && *req.NuevoVideo != "" {
		updates["video_presentacion"] = *req.NuevoVideo
	}

	if len(updates) > 0 {
		if _, err := e.DaoPostgresql.Evento.ActualizarCamposEvento(
			req.IdEvento,
			updates,
			&userID,
			&now,
		); err != nil {
			e.logger.Errorf("EditarEvento.ActualizarCamposEvento(%d): %v", req.IdEvento, err)
			return nil, &errors.InternalServerError.Default
		}
	}

	// Ubicación
	if req.NuevoLugar != nil && *req.NuevoLugar != "" {
		_, err := e.DaoPostgresql.Evento.ActualizarUbicacionEvento(
			req.IdEvento,
			*req.NuevoLugar,
			&userID,
			&now,
		)
		if err != nil {
			e.logger.Errorf("EditarEvento.ActualizarUbicacionEvento(%d): %v", req.IdEvento, err)
			return nil, &errors.InternalServerError.Default
		}
	}

	// Estado workflow (borrador/publicado/finalizado)
	if req.NuevoEstadoWorkflow != nil {
		_, err := e.DaoPostgresql.Evento.ActualizarEstadoWorkflowEvento(
			req.IdEvento,
			*req.NuevoEstadoWorkflow,
			&userID,
			&now,
		)
		if err != nil {
			e.logger.Errorf("EditarEvento.ActualizarEstadoWorkflowEvento(%d): %v", req.IdEvento, err)
			return nil, &errors.InternalServerError.Default
		}
	}

	// Estado flag (on/off)
	if req.NuevoEstadoFlag != nil {
		_, err := e.DaoPostgresql.Evento.ActualizarEstadoFlagEvento(
			req.IdEvento,
			*req.NuevoEstadoFlag,
			&userID,
			&now,
		)
		if err != nil {
			e.logger.Errorf("EditarEvento.ActualizarEstadoFlagEvento(%d): %v", req.IdEvento, err)
			return nil, &errors.InternalServerError.Default
		}
	}

	for _, f := range req.Fechas {
		// Cambiar fecha del calendario (tabla FECHA)
		if f.IdFecha != nil && f.NuevaFecha != nil && *f.NuevaFecha != "" {
			parsedDate, err := time.Parse("2006-01-02", *f.NuevaFecha)
			if err != nil {
				e.logger.Errorf("EditarEvento.Fecha parse error fecha=%s: %v", *f.NuevaFecha, err)
				return nil, &errors.UnprocessableEntityError.InvalidRequestBody
			}
			if err := e.DaoPostgresql.Evento.ActualizarFechaCalendario(
				*f.IdFecha,
				parsedDate,
				&userID,
				&now,
			); err != nil {
				e.logger.Errorf("EditarEvento.ActualizarFechaCalendario(fecha_id=%d): %v", *f.IdFecha, err)
				return nil, &errors.InternalServerError.Default
			}
		}

		// Cambiar hora inicio para un evento_fecha
		if f.IdFechaEvento > 0 && f.NuevaHoraInicio != nil && *f.NuevaHoraInicio != "" {
			parsedTime, err := time.Parse("15:04", *f.NuevaHoraInicio)
			if err != nil {
				e.logger.Errorf("EditarEvento.Hora parse error hora=%s: %v", *f.NuevaHoraInicio, err)
				return nil, &errors.UnprocessableEntityError.InvalidRequestBody
			}
			if err := e.DaoPostgresql.Evento.ActualizarHoraInicioEventoFecha(
				f.IdFechaEvento,
				parsedTime,
				&userID,
				&now,
			); err != nil {
				e.logger.Errorf("EditarEvento.ActualizarHoraInicioEventoFecha(evento_fecha_id=%d): %v", f.IdFechaEvento, err)
				return nil, &errors.InternalServerError.Default
			}
		}

		// Reasignar a otra fecha_id
		if f.IdFechaEvento > 0 && f.NuevoFechaID != nil && *f.NuevoFechaID > 0 {
			if err := e.DaoPostgresql.Evento.ReasignarFechaDeEventoFecha(
				f.IdFechaEvento,
				*f.NuevoFechaID,
				&userID,
				&now,
			); err != nil {
				e.logger.Errorf("EditarEvento.ReasignarFechaDeEventoFecha(evento_fecha_id=%d, nuevo_fecha_id=%d): %v",
					f.IdFechaEvento, *f.NuevoFechaID, err)
				return nil, &errors.InternalServerError.Default
			}
		}
	}

	for _, s := range req.Sectores {
		_, err := e.DaoPostgresql.Sector.ModificarSectorPorCampos(
			s.IdSector,
			s.SectorTipo,
			s.TotalEntradas,
			s.CantVendidas,
			s.Estado,
			&userID,
			&now,
		)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				e.logger.Warnf("EditarEvento.ModificarSectorPorCampos: sector no encontrado id=%d", s.IdSector)
				return nil, &errors.ObjectNotFoundError.EventoNotFound
			}
			e.logger.Errorf("EditarEvento.ModificarSectorPorCampos(id=%d): %v", s.IdSector, err)
			return nil, &errors.InternalServerError.Default
		}
	}

	for _, p := range req.Perfiles {
		_, err := e.DaoPostgresql.PerfilDePersona.ModificarPerfilDePersonaPorCampos(
			p.IdPerfil,
			nil, // eventoID lo dejas igual; si quieres permitir cambiarlo, mándalo en el request
			p.Nombre,
			p.Estado,
			&userID,
			&now,
		)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				e.logger.Warnf("EditarEvento.ModificarPerfilDePersona: perfil no encontrado id=%d", p.IdPerfil)
				return nil, &errors.ObjectNotFoundError.EventoNotFound
			}
			e.logger.Errorf("EditarEvento.ModificarPerfilDePersona(id=%d): %v", p.IdPerfil, err)
			return nil, &errors.InternalServerError.Default
		}
	}

	for _, ttReq := range req.TiposTicket {
		// 1) Obtener el TipoDeTicket actual de BD
		var tt model.TipoDeTicket
		if err := e.DaoPostgresql.TipoDeTicket.PostgresqlDB.
			First(&tt, "tipo_de_ticket_id = ?", ttReq.IdTipoTicket).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				e.logger.Warnf("EditarEvento.TipoDeTicket no encontrado id=%d", ttReq.IdTipoTicket)
				return nil, &errors.ObjectNotFoundError.EventoNotFound
			}
			e.logger.Errorf("EditarEvento.BuscarTipoDeTicket(id=%d): %v", ttReq.IdTipoTicket, err)
			return nil, &errors.InternalServerError.Default
		}

		// 2) Aplicar cambios parciales
		if ttReq.Nombre != nil {
			tt.Nombre = *ttReq.Nombre
		}
		if ttReq.FechaIni != nil && *ttReq.FechaIni != "" {
			if tIni, err := time.Parse("2006-01-02", *ttReq.FechaIni); err == nil {
				tt.FechaIni = tIni
			} else {
				e.logger.Errorf("EditarEvento.TipoDeTicket fechaIni parse error=%s: %v", *ttReq.FechaIni, err)
				return nil, &errors.UnprocessableEntityError.InvalidRequestBody
			}
		}
		if ttReq.FechaFin != nil && *ttReq.FechaFin != "" {
			if tFin, err := time.Parse("2006-01-02", *ttReq.FechaFin); err == nil {
				tt.FechaFin = tFin
			} else {
				e.logger.Errorf("EditarEvento.TipoDeTicket fechaFin parse error=%s: %v", *ttReq.FechaFin, err)
				return nil, &errors.UnprocessableEntityError.InvalidRequestBody
			}
		}
		if ttReq.Estado != nil {
			tt.Estado = *ttReq.Estado
		}

		// 3) Guardar
		if err := e.DaoPostgresql.TipoDeTicket.ActualizarTipoDeTicketr(&tt); err != nil {
			e.logger.Errorf("EditarEvento.ActualizarTipoDeTicketr(id=%d): %v", ttReq.IdTipoTicket, err)
			return nil, &errors.InternalServerError.Default
		}
	}

	detalle, err := e.DaoPostgresql.Evento.ObtenerEventoDetalle(req.IdEvento)
	if err != nil {
		e.logger.Errorf("EditarEvento.ObtenerEventoDetalle(idEvento=%d): %v", req.IdEvento, err)
		return nil, &errors.InternalServerError.Default
	}

	return detalle, nil
}

func (e *Evento) ObtenerTransaccionesPorEvento(eventoId string) ([]model.OrdenDeCompra, *errors.Error) {
	// Implementa la lógica para obtener las transacciones por evento
	// Esto es solo un placeholder
	var transacciones []model.OrdenDeCompra
	transacciones, err := e.DaoPostgresql.OrdenDeCompra.ListarTransaccionesPorEvento(eventoId)
	if err != nil {
		e.logger.Errorf("Failed to get transacciones for evento %s: %v", eventoId, err)
		return nil, &errors.BadRequestError.EventoNotFound
	}
	return transacciones, nil
}

func (e *Evento) PostPostgresqlInteraccionUsuarioEvento(req schemas.InteraccionConEventoRequest) (*schemas.InteraccionConEventoResponse, *errors.Error) {

	interaccionModel := model.Interaccion{
		UsuarioID:     req.UsuarioId,
		EventoID:      req.EventoId,
		Tipo:          req.Accion,
		FechaCreacion: time.Now(),
		Estado:        1,
	}

	respuesta := e.DaoPostgresql.Interaccion.CrearInteraccion(&interaccionModel)

	if respuesta != nil {
		return nil, &errors.ConflictError.InteraccionAlreadyExists
	}

	//actualizar likes
	eventoModel, _ := e.DaoPostgresql.Evento.ObtenerEventoPorId(req.EventoId)

	switch req.Accion {
	case 1:
		eventoModel.CantMeGusta++
	case 2:
		eventoModel.CantNoInteresa++
	}

	respuestaAct := e.DaoPostgresql.Evento.ActualizarInteracciones(*eventoModel)

	if respuestaAct != nil {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	eventoRes := schemas.EventoInteraccionResponse{
		EventoId:            req.EventoId,
		UsuarioId:           req.UsuarioId,
		LikesTotales:        eventoModel.CantMeGusta,
		NoMeInteresaTotales: eventoModel.CantNoInteresa,
		Like:                false,
		NoInteresa:          false,
	}

	switch req.Accion {
	case 1:
		eventoRes.Like = true
	case 2:
		eventoRes.NoInteresa = true
	}

	interaccionRes := schemas.InteraccionConEventoResponse{
		Success: true,
		Data:    eventoRes,
	}
	return &interaccionRes, nil
}

func (e *Evento) PutPostgresqlInteraccionUsuarioEvento(req schemas.InteraccionConEventoRequest) (*schemas.InteraccionConEventoResponse, *errors.Error) {

	interaccionModel := model.Interaccion{
		UsuarioID:     req.UsuarioId,
		EventoID:      req.EventoId,
		Tipo:          req.Accion,
		FechaCreacion: time.Now(),
		Estado:        1,
	}
	antInteraccion, _ := e.DaoPostgresql.Interaccion.ObtenerInteraccionesEventoUsuario(req.EventoId, req.UsuarioId)

	interaccionModel.ID = antInteraccion.ID

	respuesta := e.DaoPostgresql.Interaccion.ActualizarInteracciones(interaccionModel)

	if respuesta != nil {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	//actualizar likes
	eventoModel, _ := e.DaoPostgresql.Evento.ObtenerEventoPorId(req.EventoId)

	if req.Accion == 1 && antInteraccion.Tipo != 1 {
		eventoModel.CantMeGusta++
		if antInteraccion.Tipo == 2 {
			eventoModel.CantNoInteresa--
		}
	} else if req.Accion == 2 && antInteraccion.Tipo != 2 {
		eventoModel.CantNoInteresa++
		if antInteraccion.Tipo == 1 {
			eventoModel.CantMeGusta--
		}
	} else if req.Accion == 0 && antInteraccion.Tipo != 0 {
		if antInteraccion.Tipo == 1 {
			eventoModel.CantMeGusta--
		} else if antInteraccion.Tipo == 2 {
			eventoModel.CantNoInteresa--
		}
	}

	respuestaAct := e.DaoPostgresql.Evento.ActualizarInteracciones(*eventoModel)

	if respuestaAct != nil {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	eventoRes := schemas.EventoInteraccionResponse{
		EventoId:            req.EventoId,
		UsuarioId:           req.UsuarioId,
		LikesTotales:        eventoModel.CantMeGusta,
		NoMeInteresaTotales: eventoModel.CantNoInteresa,
		Like:                false,
		NoInteresa:          false,
	}

	switch req.Accion {
	case 1:
		eventoRes.Like = true
	case 2:
		eventoRes.NoInteresa = true
	}

	interaccionRes := schemas.InteraccionConEventoResponse{
		Success: true,
		Data:    eventoRes,
	}
	return &interaccionRes, nil
}

func (e *Evento) GetAsistentesPorEvento(eventoID int64) ([]map[string]interface{}, *errors.Error) {
	asistentes, err := e.DaoPostgresql.Evento.ObtenerAsistentesPorEvento(eventoID)
	if err != nil {
		e.logger.Errorf("❌ [ADAPTER] Error obteniendo asistentes: %v", err)
		return nil, &errors.InternalServerError.Default
	}

	if len(asistentes) == 0 {
		e.logger.Infof("ℹ️ [ADAPTER] No se encontraron asistentes para el evento: %d", eventoID)
		return []map[string]interface{}{}, nil
	}

	return asistentes, nil
}
