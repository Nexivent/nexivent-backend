package repository

import (
	"time"

	"github.com/Loui27/nexivent-backend/internal/dao/model"
	"github.com/Loui27/nexivent-backend/internal/schemas"
	"github.com/Loui27/nexivent-backend/logging"
	"gorm.io/gorm"
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

func (e *Evento) ObtenerEventosDisponiblesSinFiltros(limit int, page int) (*schemas.EventosPaginados, error) {
	var categoriaID *int64
	var titulo *string
	var descripcion *string
	var lugar *string
	var fecha *time.Time
	var horaInicio *time.Time

	eventos, respuesta := e.ObtenerEventosDisponiblesConFiltros(
		categoriaID, titulo, descripcion, lugar, fecha, horaInicio, limit, page,
	)
	if respuesta != nil {
		return nil, respuesta
	}

	return eventos, nil
}

func (e *Evento) ObtenerEventosDisponiblesConFiltros(
	categoriaID *int64,
	titulo *string,
	descripcion *string,
	lugar *string,
	fecha *time.Time,
	horaInicio *time.Time,
	limit int,
	page int,
) (*schemas.EventosPaginados, error) {
	var (
		eventos []*model.Evento
		total   int64
	)

	// Calcular offset
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	// Construcción base del query
	query := e.PostgresqlDB.
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
	queryCount := query.Session(&gorm.Session{}) // Clona el query sin afectar el original
	queryCount.Count(&total)

	// Aplicar paginación
	respuesta := query.
		Order("f.fecha_evento ASC").
		Limit(limit).
		Offset(offset).
		Find(&eventos)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	// Calcular total de páginas
	totalPaginas := int((total + int64(limit) - 1) / int64(limit))

	// Retornar resultado completo
	resultado := &schemas.EventosPaginados{
		Eventos:      eventos,
		Total:        total,
		PaginaActual: page,
		TotalPaginas: totalPaginas,
	}

	return resultado, nil
}
