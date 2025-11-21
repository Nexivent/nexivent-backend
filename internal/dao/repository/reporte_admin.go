package repository

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"gorm.io/gorm"
)

// GenerarReporteAdmin ejecuta las 4 consultas en paralelo (conceptualmente) reutilizando filtros
func (e *Evento) GenerarReporteAdmin(
	fechaInicio, fechaFin *time.Time,
	idCategoria, idOrganizador *int64,
	estadoInt *int16, // Estado convertido a entero (ej. 1=Publicado)
	limit int,
) (*schemas.AdminReportResponse, error) {

	var response schemas.AdminReportResponse

	// 1. Definir el Scope de filtros (Reutilizable)
	filtros := func(db *gorm.DB) *gorm.DB {
		// Joins necesarios para filtrar por fecha o categoría
		query := db.Joins("JOIN evento_fecha ef ON ef.evento_id = evento.evento_id").
			Joins("LEFT JOIN categoria c ON c.categoria_id = evento.categoria_id")

		if fechaInicio != nil {
			query = query.Where("ef.hora_inicio >= ?", *fechaInicio)
		}
		if fechaFin != nil {
			query = query.Where("ef.hora_inicio <= ?", *fechaFin)
		}
		if idCategoria != nil {
			query = query.Where("evento.categoria_id = ?", *idCategoria)
		}
		if idOrganizador != nil {
			query = query.Where("evento.organizador_id = ?", *idOrganizador)
		}
		if estadoInt != nil {
			query = query.Where("evento.evento_estado = ?", *estadoInt) // O el campo de estado que uses
		}
		return query
	}

	// --- CONSULTA A: SUMMARY (Métricas Generales) ---
	// Usamos CASE WHEN para contar estados en una sola pasada
	err := e.PostgresqlDB.Table("evento").
		Scopes(filtros).
		Select(`
			COUNT(DISTINCT evento.evento_id) as total_eventos,
			COALESCE(SUM(CASE WHEN evento.evento_estado = 1 THEN 1 ELSE 0 END), 0) as total_publicados,
			COALESCE(SUM(CASE WHEN evento.evento_estado = 2 THEN 1 ELSE 0 END), 0) as total_cancelados,
			COALESCE(SUM(CASE WHEN evento.evento_estado = 0 THEN 1 ELSE 0 END), 0) as total_borradores,
			COALESCE(SUM(evento.cant_vendido_total), 0) as entradas_vendidas_totales,
			COALESCE(SUM(evento.total_recaudado), 0) as recaudacion_total
		`).
		Scan(&response.Summary).Error

	if err != nil {
		return nil, err
	}

	// Si no hay eventos, retornamos nil para que el Adapter maneje el 204/404
	if response.Summary.TotalEventos == 0 {
		return nil, nil
	}

	// --- CONSULTA B: LISTADO DETALLADO ---
	err = e.PostgresqlDB.Table("evento").
		Scopes(filtros).
		Select(`
			DISTINCT ON (evento.evento_id) evento.evento_id,
			evento.titulo,
			c.nombre as categoria,
			evento.lugar,
			evento.evento_estado, -- Se mapeará en el Adapter
			ef.hora_inicio as fecha_inicio,
			evento.cant_vendido_total as entradas_vendidas,
			evento.total_recaudado as recaudacion_total
		`).
		Limit(limit).
		Scan(&response.Events).Error
	if err != nil {
		return nil, err
	}

	// --- CONSULTA C: TOP EVENTOS (Por Recaudación) ---
	err = e.PostgresqlDB.Table("evento").
		Scopes(filtros).
		Select(`
			DISTINCT ON (evento.evento_id) evento.evento_id,
			evento.titulo,
			evento.lugar,
			evento.cant_vendido_total as entradas_vendidas,
			evento.total_recaudado as recaudacion
		`).
		Order("evento.total_recaudado DESC").
		Limit(5). // Top 5 fijo
		Scan(&response.TopEventos).Error
	if err != nil {
		return nil, err
	}

	// --- CONSULTA D: AGRUPADO POR CATEGORÍA ---
	err = e.PostgresqlDB.Table("evento").
		Scopes(filtros).
		Select(`
			c.categoria_id as id_categoria,
			c.nombre as categoria,
			COUNT(DISTINCT evento.evento_id) as cantidad_eventos,
			COALESCE(SUM(evento.total_recaudado), 0) as recaudacion_total,
			COALESCE(SUM(evento.cant_vendido_total), 0) as entradas_vendidas
		`).
		Group("c.categoria_id, c.nombre").
		Scan(&response.ByCategory).Error
	if err != nil {
		return nil, err
	}

	return &response, nil
}
