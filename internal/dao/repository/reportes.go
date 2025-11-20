package repository

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
)

type Reportes struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewReportesController(
	logger logging.Logger,
	postgresqlDB *gorm.DB,
) *OrdenDeCompra {
	return &OrdenDeCompra{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

func (r *Reportes) BuscarEventosParaReporte(
	fechaDesde, fechaHasta *time.Time,
	idEvento *int64,
) ([]model.Evento, error) {

	query := r.PostgresqlDB.Model(&model.Evento{}).
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
	if err := query.Group("evento.evento_id").Find(&eventos).Error; err != nil {
		return nil, err
	}

	return eventos, nil
}
