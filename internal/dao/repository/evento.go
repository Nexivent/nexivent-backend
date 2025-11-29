package repository

import (
	"strings"
	"time"

	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Evento struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewEventoController(
	logger logging.Logger,
	postgresqlDB *gorm.DB,
) *Evento {
	return &Evento{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

func (e *Evento) CrearEvento(Evento *model.Evento) error {
	respuesta := e.PostgresqlDB.Create(Evento)
	if respuesta.Error != nil {
		return respuesta.Error
	}
	return nil
}

func (e *Evento) ObtenerEventosDisponiblesSinFiltros() ([]*model.Evento, error) {
	var categoriaID *int64
	var titulo *string
	var descripcion *string
	var lugar *string
	var fecha *time.Time
	var horaInicio *time.Time
	var estado *int16
	soloFuturos := false

	eventos, respuesta := e.ObtenerEventosDisponiblesConFiltros(
		categoriaID, nil, titulo, descripcion, lugar, fecha, horaInicio, estado, soloFuturos,
	)
	if respuesta != nil {
		return nil, respuesta
	}

	return eventos, nil
}

func (e *Evento) ObtenerEventosDisponiblesConFiltros(
	categoriaID *int64,
	organizadorID *int64,
	titulo *string,
	descripcion *string,
	lugar *string,
	fecha *time.Time,
	horaInicio *time.Time,
	estado *int16,
	soloFuturos bool,
) ([]*model.Evento, error) {
	var eventos []*model.Evento

	// Construcción base del query
	query := e.PostgresqlDB.
		Model(&model.Evento{}).
		Select("DISTINCT ON (evento.evento_id) evento.*").
		Preload("Fechas.Fecha").
		Preload("Sectores").
		Preload("Sectores.Tarifa").
		Preload("Sectores.Tarifa.TipoDeTicket").
		Preload("Sectores.Tarifa.PerfilPersona").
		Preload("TiposTicket").
		Joins("JOIN evento_fecha ef ON ef.evento_id = evento.evento_id").
		Joins("JOIN fecha f ON f.fecha_id = ef.fecha_id")

	// Aplicar filtros dinámicamente
	if categoriaID != nil {
		query = query.Where("evento.categoria_id = ?", *categoriaID)
	}
	if fecha != nil {
		query = query.Where("f.fecha_evento = ?", *fecha)
	}
	if horaInicio != nil {
		query = query.Where("ef.hora_inicio = ?", *horaInicio)
	}
	if soloFuturos {
		hoy := time.Now().Truncate(24 * time.Hour)
		query = query.Where("f.fecha_evento >= ?", hoy)
	}
	if organizadorID != nil {
		query = query.Where("evento.organizador_id = ?", *organizadorID)
	}
	if estado != nil {
		query = query.Where("evento.evento_estado = ?", *estado)
	}

	// Filtro OR agrupado para búsqueda textual
	var condiciones []string
	var valores []interface{}

	if titulo != nil && *titulo != "" {
		condiciones = append(condiciones, "evento.titulo ILIKE ?")
		valores = append(valores, "%"+*titulo+"%")
	}

	if descripcion != nil && *descripcion != "" {
		condiciones = append(condiciones, "evento.descripcion ILIKE ?")
		valores = append(valores, "%"+*descripcion+"%")
	}

	if lugar != nil && *lugar != "" {
		condiciones = append(condiciones, "evento.lugar ILIKE ?")
		valores = append(valores, "%"+*lugar+"%")
	}

	// Solo agregar el OR si al menos un campo se envió
	if len(condiciones) > 0 {
		orGroup := "(" + strings.Join(condiciones, " OR ") + ")"
		query = query.Where(orGroup, valores...)
	}

	// Aplicar paginación
	respuesta := query.
		Order("evento.evento_id, f.fecha_evento ASC").
		//Limit(limit).
		//Offset(offset).
		Find(&eventos)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	return eventos, nil
}

func (e *Evento) ObtenerEventosParaElFeed(usuarioId *int64) ([]*model.Evento, error) {
	// Construcción base del query
	var eventos []*model.Evento
	query := e.PostgresqlDB.
		Preload("Fechas").
		Preload("Fechas.Fecha").
		Preload("Sectores").
		Preload("Sectores.Tarifa").
		Preload("Sectores.Tarifa.TipoDeTicket").
		Preload("Sectores.Tarifa.PerfilPersona").
		Preload("TiposTicket").
		Joins("JOIN evento_fecha ef ON ef.evento_id = evento.evento_id").
		Joins("JOIN fecha f ON f.fecha_id = ef.fecha_id").
		Where("f.fecha_evento >= CURRENT_DATE").
		Where("evento.evento_estado = 1").
		Where("evento.estado = 1").
		Where(`f.fecha_evento = ( SELECT MIN(f2.fecha_evento)
        	FROM fecha f2
        	JOIN evento_fecha ef2 ON ef2.fecha_id = f2.fecha_id
        	WHERE ef2.evento_id = evento.evento_id)`)

	if usuarioId != nil {
		query = query.Where(`NOT EXISTS (SELECT 1 FROM interaccion i WHERE i.usuario_id = ? AND i.evento_id = evento.evento_id)`, usuarioId)
	}

	respuesta := query.
		Order("((2*evento.cant_me_gusta - evento.cant_no_interesa) / GREATEST(1, (f.fecha_evento::date - CURRENT_DATE))) DESC").
		Find(&eventos)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	return eventos, nil
}

func (e *Evento) CargarEventosNuevamenteParaElFeed(usuarioId *int64) ([]*model.Evento, error) {
	var eventos []*model.Evento

	// Query base
	query := e.PostgresqlDB.
		Preload("Fechas").
		Preload("Fechas.Fecha").
		Preload("Sectores").
		Preload("Sectores.Tarifa").
		Preload("Sectores.Tarifa.TipoDeTicket").
		Preload("Sectores.Tarifa.PerfilPersona").
		Preload("TiposTicket").
		Joins("JOIN evento_fecha ef ON ef.evento_id = evento.evento_id").
		Joins("JOIN fecha f ON f.fecha_id = ef.fecha_id").
		Where("f.fecha_evento >= CURRENT_DATE").
		Where("evento.evento_estado = 1").
		Where("evento.estado = 1").
		Where(`f.fecha_evento = (
            SELECT MIN(f2.fecha_evento)
            FROM fecha f2
            JOIN evento_fecha ef2 ON ef2.fecha_id = f2.fecha_id
            WHERE ef2.evento_id = evento.evento_id
        )`)
	// solo pre-cargamos la interacción del usuario (si existe)
	if usuarioId != nil {
		query = query.Preload("Interacciones", "usuario_id = ?", *usuarioId)
	}

	respuesta := query.
		Order("((2*evento.cant_me_gusta - evento.cant_no_interesa) / GREATEST(1, (f.fecha_evento::date - CURRENT_DATE))) DESC").
		Find(&eventos)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	return eventos, nil
}

// ===============================
//
//	Actualización: UBICACIÓN
//
// ===============================
func (e *Evento) ActualizarUbicacionEvento(
	eventoID int64,
	nuevoLugar string,
	usuarioModificacion *int64,
	fechaModificacion *time.Time,
) (*model.Evento, error) {

	if eventoID <= 0 || nuevoLugar == "" {
		return nil, gorm.ErrInvalidData
	}

	updates := map[string]any{
		"lugar": nuevoLugar,
	}
	if usuarioModificacion != nil {
		updates["usuario_modificacion"] = *usuarioModificacion
	}
	if fechaModificacion != nil {
		updates["fecha_modificacion"] = *fechaModificacion
	}

	var ev model.Evento
	res := e.PostgresqlDB.
		Model(&ev).
		Clauses(clause.Returning{}).
		Where("evento_id = ?", eventoID).
		Updates(updates)

	if res.Error != nil {
		e.logger.Errorf("ActualizarUbicacionEvento id=%d: %v", eventoID, res.Error)
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &ev, nil
}

// =======================================
//
//	Actualización: ESTADO WORKFLOW (evento_estado)
//	(p.ej., 0=Borrador, 1=Publicado, 2=Finalizado)
//
// =======================================
func (e *Evento) ActualizarEstadoWorkflowEvento(
	eventoID int64,
	nuevoEstado int16,
	usuarioModificacion *int64,
	fechaModificacion *time.Time,
) (*model.Evento, error) {

	if eventoID <= 0 {
		return nil, gorm.ErrInvalidData
	}

	updates := map[string]any{
		"evento_estado": nuevoEstado,
	}
	if usuarioModificacion != nil {
		updates["usuario_modificacion"] = *usuarioModificacion
	}
	if fechaModificacion != nil {
		updates["fecha_modificacion"] = *fechaModificacion
	}

	var ev model.Evento
	res := e.PostgresqlDB.
		Model(&ev).
		Clauses(clause.Returning{}).
		Where("evento_id = ?", eventoID).
		Updates(updates)

	if res.Error != nil {
		e.logger.Errorf("ActualizarEstadoWorkflowEvento id=%d: %v", eventoID, res.Error)
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &ev, nil
}

// =======================================
//
//	Actualización: ESTADO FLAG (estado)
//	(0=inactivo, 1=activo; soft on/off del registro)
//
// =======================================
func (e *Evento) ActualizarEstadoFlagEvento(
	eventoID int64,
	nuevoEstado int16,
	usuarioModificacion *int64,
	fechaModificacion *time.Time,
) (*model.Evento, error) {

	if eventoID <= 0 {
		return nil, gorm.ErrInvalidData
	}

	updates := map[string]any{
		"estado": nuevoEstado,
	}
	if usuarioModificacion != nil {
		updates["usuario_modificacion"] = *usuarioModificacion
	}
	if fechaModificacion != nil {
		updates["fecha_modificacion"] = *fechaModificacion
	}

	var ev model.Evento
	res := e.PostgresqlDB.
		Model(&ev).
		Clauses(clause.Returning{}).
		Where("evento_id = ?", eventoID).
		Updates(updates)

	if res.Error != nil {
		e.logger.Errorf("ActualizarEstadoFlagEvento id=%d: %v", eventoID, res.Error)
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &ev, nil
}

// ActualizarCamposEvento actualiza columnas puntuales del evento en una sola llamada.
func (e *Evento) ActualizarCamposEvento(
	eventoID int64,
	updates map[string]any,
	usuarioModificacion *int64,
	fechaModificacion *time.Time,
) (*model.Evento, error) {
	if eventoID <= 0 || len(updates) == 0 {
		return nil, gorm.ErrInvalidData
	}

	if usuarioModificacion != nil {
		updates["usuario_modificacion"] = *usuarioModificacion
	}
	if fechaModificacion != nil {
		updates["fecha_modificacion"] = *fechaModificacion
	}

	var ev model.Evento
	res := e.PostgresqlDB.
		Model(&ev).
		Clauses(clause.Returning{}).
		Where("evento_id = ?", eventoID).
		Updates(updates)

	if res.Error != nil {
		e.logger.Errorf("ActualizarCamposEvento id=%d: %v", eventoID, res.Error)
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &ev, nil
}

// =====================================================
//  FECHAS (tablas: fecha, evento_fecha)
//  - Cambiar fecha del calendario (tabla fecha)
//  - Cambiar hora inicio de un evento_fecha
//  - Reasignar fecha en un evento_fecha (cambiar fecha_id)
// =====================================================

// Cambia el valor de fecha_evento (tabla FECHA) para un fecha_id dado.
// Ojo: este cambio afecta a todos los evento_fecha que referencien ese fecha_id.
func (e *Evento) ActualizarFechaCalendario(
	fechaID int64,
	nuevaFecha time.Time, // usar solo la parte de día acorde a tu diseño
	usuarioModificacion *int64,
	fechaModificacion *time.Time,
) error {

	if fechaID <= 0 {
		return gorm.ErrInvalidData
	}

	updates := map[string]any{
		"fecha_evento": nuevaFecha, // Postgres DATE (GORM hace el cast si tu model.Fecha es time.Time)
	}
	// La tabla fecha no tiene campos de auditoría en tu DDL; si luego los agregas, setéalos aquí.
	if usuarioModificacion != nil {
		updates["usuario_modificacion"] = *usuarioModificacion
	}
	if fechaModificacion != nil {
		updates["fecha_modificacion"] = *fechaModificacion
	}

	res := e.PostgresqlDB.
		Table("fecha").
		Where("fecha_id = ?", fechaID).
		Updates(updates)

	if res.Error != nil {
		e.logger.Errorf("ActualizarFechaCalendario fecha_id=%d: %v", fechaID, res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Cambia la HORA de inicio de un registro evento_fecha (no la fecha).
func (e *Evento) ActualizarHoraInicioEventoFecha(
	eventoFechaID int64,
	nuevaHora time.Time, // usa time con la hora deseada (Postgres TIMESTAMPTZ)
	usuarioModificacion *int64,
	fechaModificacion *time.Time,
) error {

	if eventoFechaID <= 0 {
		return gorm.ErrInvalidData
	}

	updates := map[string]any{
		"hora_inicio": nuevaHora,
	}
	if usuarioModificacion != nil {
		updates["usuario_modificacion"] = *usuarioModificacion
	}
	if fechaModificacion != nil {
		updates["fecha_modificacion"] = *fechaModificacion
	}

	res := e.PostgresqlDB.
		Table("evento_fecha").
		Where("evento_fecha_id = ?", eventoFechaID).
		Updates(updates)

	if res.Error != nil {
		e.logger.Errorf("ActualizarHoraInicioEventoFecha evento_fecha_id=%d: %v", eventoFechaID, res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Reasigna la fecha (fecha_id) de un evento_fecha específico.
// Útil si creas una nueva fecha en 'fecha' y quieres apuntar el evento_fecha a esa nueva fecha.
func (e *Evento) ReasignarFechaDeEventoFecha(
	eventoFechaID int64,
	nuevoFechaID int64,
	usuarioModificacion *int64,
	fechaModificacion *time.Time,
) error {

	if eventoFechaID <= 0 || nuevoFechaID <= 0 {
		return gorm.ErrInvalidData
	}

	updates := map[string]any{
		"fecha_id": nuevoFechaID,
	}
	if usuarioModificacion != nil {
		updates["usuario_modificacion"] = *usuarioModificacion
	}
	if fechaModificacion != nil {
		updates["fecha_modificacion"] = *fechaModificacion
	}

	res := e.PostgresqlDB.
		Table("evento_fecha").
		Where("evento_fecha_id = ?", eventoFechaID).
		Updates(updates)

	if res.Error != nil {
		e.logger.Errorf("ReasignarFechaDeEventoFecha evento_fecha_id=%d: %v", eventoFechaID, res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

/*
func (e *Evento) BuscarEventosParaReporte(
	fechaDesde, fechaHasta *time.Time,
	idEvento *int64,
) ([]model.Evento, error) {

	query := e.PostgresqlDB.Model(&model.Evento{}).
		Joins("JOIN evento_fecha ef ON ef.evento_id = evento.evento_id").
		Joins("JOIN fecha f ON f.fecha_id = ef.fecha_id")

	if idEvento != nil {
		query = query.Where("evento.id = ?", *idEvento)
	}
	if fechaDesde != nil {
		query = query.Where("f.fecha_evento >= ?", *fechaDesde)
	}
	if fechaHasta != nil {
		query = query.Where("f.fecha_evento <= ?", *fechaHasta)
	}

	var eventos []model.Evento

	if err := query.Group("evento.evento_id").
		Preload("Fechas").
		Preload("Fechas.Fecha").
		Preload("TipoDeTicket").
		Preload("TipoDeTicket.Tarifa").
		Find(&eventos).Error; err != nil {
		return nil, err
	}

	return eventos, nil
}*/

func (e *Evento) ObtenerEventosDelOrganizador(idOrganizador int64) ([]*model.Evento, error) {
	var eventos []*model.Evento

	res := e.PostgresqlDB.Table("evento").
		Preload("Fechas").
		Preload("Fechas.Fecha").
		Where("organizador_id = ? AND evento_estado=1", idOrganizador).
		Find(&eventos)

	if res.Error != nil {
		return nil, res.Error
	}

	return eventos, nil
}

func (e *Evento) ObtenerEventoPorId(id int64) (*model.Evento, error) {
	var evento *model.Evento

	res := e.PostgresqlDB.Table("evento").
		Preload("Fechas").
		Preload("Fechas.Fecha").
		Where("evento_id = ? AND evento_estado = 1", id).
		Find(&evento)

	if res.Error != nil {
		return nil, res.Error
	}

	return evento, nil
}

func (e *Evento) ObtenerEventoDetalle(eventoId int64) (*schemas.EventoDetalleDTO, error) {
	// Obtener datos básicos del evento
	var eventoBase model.Evento

	respuesta := e.PostgresqlDB.
		Table("evento").
		Select("evento_id , titulo, descripcion, imagen_portada, lugar").
		Where("evento_id = ?", eventoId).
		First(&eventoBase)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	// Obtener fechas del evento (join evento_fecha con fecha)
	//var fechas []FechaEventoDTO
	var fechas []schemas.FechaEventoDTO
	e.PostgresqlDB.
		Table("evento_fecha ef").
		Select(`ef.evento_fecha_id as id_fecha_evento,
        TO_CHAR(f.fecha_evento, 'YYYY-MM-DD') as fecha,
        TO_CHAR(ef.hora_inicio, 'HH24:MI') as hora_inicio,
        '' as hora_fin,
        ef.ganancia_neta_organizador as ganancia_neta_organizador`).
		Joins("JOIN fecha f ON ef.fecha_id = f.fecha_id").
		Where("ef.evento_id = ? AND ef.estado = 1", eventoId).
		Find(&fechas)

	// Obtener tarifas con joins
	//var tarifas []TarifaDTO
	/*var tarifas []schemas.TarifaDTO
	e.PostgresqlDB.
		Table("tarifa t").
		Select(`t.tarifa_id as id_tarifa,
			t.precio,
			s.nombre_sector as tipo_sector,
			(s.stock - s.cant_vendidas) as stock_disponible,
			tt.nombre_ticket as tipo_ticket,
			TO_CHAR(t.fecha_ini, 'YYYY-MM-DD') as fecha_ini,
			TO_CHAR(t.fecha_fin, 'YYYY-MM-DD') as fecha_fin,
			pp.perfil`).
		Joins("JOIN sectores s ON t.sector_id = s.sector_id").
		Joins("JOIN tipos_ticket tt ON t.tipo_ticket_id = tt.tipo_ticket_id").
		Joins("JOIN perfil_persona pp ON t.perfil_persona_id = pp.perfil_persona_id").
		Where("s.evento_id = ?", eventoId).
		Find(&tarifas)
	*/
	var tarifas []schemas.TarifaDTO
	e.PostgresqlDB.
		Table("tarifa t").
		Select(`t.tarifa_id as id_tarifa,
			t.precio, s.sector_id as id_sector,
			s.sector_tipo as tipo_sector,
			(s.total_entradas - s.cant_vendidas) as stock_disponible, tt.tipo_de_ticket_id as id_tipo_ticket,
			tt.nombre as tipo_ticket,
			TO_CHAR(tt.fecha_ini, 'YYYY-MM-DD') as fecha_ini,
			TO_CHAR(tt.fecha_fin, 'YYYY-MM-DD') as fecha_fin, pp.perfil_de_persona_id as id_perfil,
			COALESCE(pp.nombre, '') as perfil`).
		Joins("JOIN sector s ON t.sector_id = s.sector_id").
		Joins("JOIN tipo_de_ticket tt ON t.tipo_de_ticket_id = tt.tipo_de_ticket_id").
		Joins("LEFT JOIN perfil_de_persona pp ON t.perfil_de_persona_id = pp.perfil_de_persona_id").
		Where("s.evento_id = ? AND t.estado = 1", eventoId).
		Find(&tarifas)

	return &schemas.EventoDetalleDTO{
		IDEvento:      eventoBase.ID,
		Titulo:        eventoBase.Titulo,
		Descripcion:   eventoBase.Descripcion,
		Lugar:         eventoBase.Lugar,
		ImagenPortada: eventoBase.ImagenPortada,
		Fechas:        fechas,
		Tarifas:       tarifas,
	}, nil
}

func (e *Evento) ActualizarInteracciones(evento model.Evento) error {
	respuesta := e.PostgresqlDB.
		Table("evento").
		Where("evento_id = ?", evento.ID).
		Update("cant_me_gusta", evento.CantMeGusta).
		Update("cant_no_interesa", evento.CantNoInteresa)

	if respuesta != nil {
		return respuesta.Error
	}
	return nil
}

func (e *Evento) ObtenerAsistentesPorEvento(eventoID int64) ([]map[string]interface{}, error) {
	e.logger.Infof("📋 [REPO] Obteniendo asistentes del evento ID: %d", eventoID)

	var asistentes []map[string]interface{}

	query := `
        SELECT DISTINCT
            u.usuario_id as id,
            u.correo as email,
            u.nombre as nombre,
            COUNT(DISTINCT t.ticket_id) as cantidad_tickets,
            SUM(DISTINCT odc.total) as total_gastado
        FROM usuario u
        INNER JOIN orden_de_compra odc ON u.usuario_id = odc.usuario_id
        INNER JOIN ticket t ON t.orden_de_compra_id = odc.orden_de_compra_id
        INNER JOIN evento_fecha ef ON t.evento_fecha_id = ef.evento_fecha_id
        INNER JOIN evento ev ON ef.evento_id = ev.evento_id
        WHERE ev.evento_id = $1
          AND odc.estado_de_orden = 1
          AND t.estado_de_ticket = 1
          AND u.usuario_id != ev.organizador_id
        GROUP BY u.usuario_id, u.correo, u.nombre
        ORDER BY u.nombre ASC
    `

	rows, err := e.PostgresqlDB.Raw(query, eventoID).Rows()
	if err != nil {
		e.logger.Errorf("❌ [REPO] Error ejecutando query de asistentes: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var asistente struct {
			ID              int64
			Email           string
			Nombre          string
			CantidadTickets int
			TotalGastado    float64
		}

		if err := rows.Scan(
			&asistente.ID,
			&asistente.Email,
			&asistente.Nombre,
			&asistente.CantidadTickets,
			&asistente.TotalGastado,
		); err != nil {
			e.logger.Errorf("❌ [REPO] Error escaneando asistente: %v", err)
			continue
		}

		asistentes = append(asistentes, map[string]interface{}{
			"id":               asistente.ID,
			"email":            asistente.Email,
			"nombre":           asistente.Nombre,
			"cantidad_tickets": asistente.CantidadTickets,
			"total_gastado":    asistente.TotalGastado,
		})
	}

	e.logger.Infof("✅ [REPO] Asistentes encontrados: %d", len(asistentes))
	return asistentes, nil
}
