package repository

import (
	"github.com/Loui27/nexivent-backend/internal/dao/model"
	"github.com/Loui27/nexivent-backend/logging"
	"gorm.io/gorm"
	// "github.com/Loui27/nexivent-backend/internal/dao/model"
)

type EventoFecha struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewEventoFechaController(
	logger logging.Logger,
	postgreesqlDB *gorm.DB,
) *EventoFecha {
	return &EventoFecha{
		logger:       logger,
		PostgresqlDB: postgreesqlDB,
	}
}

func (e *EventoFecha) CrearEventoFecha(EventoFecha *model.EventoFecha) error {
	respuesta := e.PostgresqlDB.Create(EventoFecha)

	if respuesta.Error != nil {
		return respuesta.Error
	}

	return nil
}
