package repository

import (
	"github.com/Loui27/nexivent-backend/internal/dao/model"
	"github.com/Loui27/nexivent-backend/logging"
	"gorm.io/gorm"
	// "github.com/Loui27/nexivent-backend/internal/dao/model"
)

type Tarifa struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewTarifaController(
	logger logging.Logger,
	postgreesqlDB *gorm.DB,
) *Tarifa {
	return &Tarifa{
		logger:       logger,
		PostgresqlDB: postgreesqlDB,
	}
}

func (t *Tarifa) CrearTarifa(Tarifa *model.Tarifa) error {
	resultado := t.PostgresqlDB.Create(Tarifa)

	if resultado.Error != nil {
		return resultado.Error
	}

	return nil
}

func (t *Tarifa) ActualizarTarifa(Tarifa *model.Tarifa) error {
	respuesta := t.PostgresqlDB.Save(Tarifa)

	if respuesta.Error != nil {
		return respuesta.Error
	}

	return nil
}
