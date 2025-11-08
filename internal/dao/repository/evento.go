package repository

import (
	"github.com/Loui27/nexivent-backend/logging"
	"gorm.io/gorm"
	// "github.com/Loui27/nexivent-backend/internal/dao/model"
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

func (e *Evento) CrearEvento(Evento *model.Evento) (int64, error) {
	respuesta := e.PostgresqlDB.Create(Evento)
	if respuesta.Error != nil{
		return 0, respuesta.Error
	}

	return respuesta.ID, nil
}