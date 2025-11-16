package repository

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
	// "github.com/Loui27/nexivent-backend/internal/dao/model"
)

type Fecha struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewFechaController(
	logger logging.Logger,
	postgreesqlDB *gorm.DB,
) *Fecha {
	return &Fecha{
		logger:       logger,
		PostgresqlDB: postgreesqlDB,
	}
}

func (f *Fecha) CrearFecha(Fecha *model.Fecha) error {
	respuesta := f.PostgresqlDB.Create(Fecha)

	if respuesta.Error != nil {
		return respuesta.Error
	}

	return nil
}

func (f *Fecha) ObtenerPorFecha(Fecha time.Time) (*model.Fecha, error) {
	var fechaObtenida *model.Fecha
	respuesta := f.PostgresqlDB.Where("fecha_evento = ?", Fecha).Find(&fechaObtenida)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	return fechaObtenida, nil
}
