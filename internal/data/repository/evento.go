package repository

import (
	"fmt"
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventoSQL struct {
	DB *gorm.DB
}

func (e *EventoSQL) CrearEvento(Evento *model.Evento) error {
	respuesta := e.DB.Create(Evento)
	if respuesta.Error != nil {
		return respuesta.Error
	}
	fmt.Printf("Evento creado: %+v\n", Evento)
	return nil
}

func (e *EventoSQL) ObtenerEventosDisponiblesSinFiltros() ([]*model.Evento, error) {
	var categoriaID *int64
	var titulo *string
	var descripcion *string
	var lugar *string
	var fecha *time.Time
	var horaInicio *time.Time

	eventos, respuesta := e.ObtenerEventosDisponiblesConFiltros(
		categoriaID, titulo, descripcion, lugar, fecha, horaInicio,
	)
	if respuesta != nil {
		return nil, respuesta
	}

	return eventos, nil
}

func (e *EventoSQL) ObtenerEventosDisponiblesConFiltros(
	categoriaID *int64,
	titulo *string,
	descripcion *string,
	lugar *string,
	fecha *time.Time,
	horaInicio *time.Time,
) ([]*model.Evento, error) {
	var eventos []*model.Evento

	// Construcción base del query
	query := e.DB.
		Joins("JOIN evento_fecha ef ON ef.evento_id = evento.evento_id").
		Joins("JOIN fecha f ON f.fecha_id = ef.fecha_id").
		Where("evento.evento_estado = 1 AND f.fecha_evento > NOW()")

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
	if titulo != nil && *titulo != "" {
		query = query.Where("evento.titulo ILIKE ?", "%"+*titulo+"%")
	}
	if descripcion != nil && *descripcion != "" {
		query = query.Where("evento.descripcion ILIKE ?", "%"+*descripcion+"%")
	}
	if lugar != nil && *lugar != "" {
		query = query.Where("evento.lugar ILIKE ?", "%"+*lugar+"%")
	}

	// Contar total antes de aplicar limit/offset
	//queryCount := query.Session(&gorm.Session{}) // Clona el query sin afectar el original
	//queryCount.Count(&total)

	// Aplicar paginación
	respuesta := query.
		Order("f.fecha_evento ASC").
		//Limit(limit).
		//Offset(offset).
		Find(&eventos)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	// Calcular total de páginas
	//totalPaginas := int((total + int64(limit) - 1) / int64(limit))

	// Retornar resultado completo
	/*resultado := &schemas.EventosPaginados{
		Eventos:      eventos,
		Total:        total,
		PaginaActual: page,
		TotalPaginas: totalPaginas,
	}*/

	return eventos, nil
}

// ===============================
//
//	Actualización: UBICACIÓN
//
// ===============================
func (e *EventoSQL) ActualizarUbicacionEvento(
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
	res := e.DB.
		Model(&ev).
		Clauses(clause.Returning{}).
		Where("evento_id = ?", eventoID).
		Updates(updates)

	if res.Error != nil {
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
func (e *EventoSQL) ActualizarEstadoWorkflowEvento(
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
	res := e.DB.
		Model(&ev).
		Clauses(clause.Returning{}).
		Where("evento_id = ?", eventoID).
		Updates(updates)

	if res.Error != nil {
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
func (e *EventoSQL) ActualizarEstadoFlagEvento(
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
	res := e.DB.
		Model(&ev).
		Clauses(clause.Returning{}).
		Where("evento_id = ?", eventoID).
		Updates(updates)

	if res.Error != nil {
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
func (e *EventoSQL) ActualizarFechaCalendario(
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

	res := e.DB.
		Table("fecha").
		Where("fecha_id = ?", fechaID).
		Updates(updates)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Cambia la HORA de inicio de un registro evento_fecha (no la fecha).
func (e *EventoSQL) ActualizarHoraInicioEventoFecha(
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

	res := e.DB.
		Table("evento_fecha").
		Where("evento_fecha_id = ?", eventoFechaID).
		Updates(updates)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Reasigna la fecha (fecha_id) de un evento_fecha específico.
// Útil si creas una nueva fecha en 'fecha' y quieres apuntar el evento_fecha a esa nueva fecha.
func (e *EventoSQL) ReasignarFechaDeEventoFecha(
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

	res := e.DB.
		Table("evento_fecha").
		Where("evento_fecha_id = ?", eventoFechaID).
		Updates(updates)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
