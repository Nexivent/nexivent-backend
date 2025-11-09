package repository

import (
	"github.com/Loui27/nexivent-backend/internal/dao/model"
	"github.com/Loui27/nexivent-backend/logging"
	"gorm.io/gorm"
)

type Sector struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewSectorController(
	logger logging.Logger,
	postgreesqlDB *gorm.DB,
) *Sector {
	return &Sector{
		logger:       logger,
		PostgresqlDB: postgreesqlDB,
	}
}

func (s *Sector) CrearSector(sector *model.Sector) error {
	resultado := s.PostgresqlDB.Create(sector)

	if resultado.Error != nil {
		return resultado.Error
	}

	return nil
}

func (s *Sector) ActualizarSector(sector *model.Sector) error {
	respuesta := s.PostgresqlDB.Save(sector)

	if respuesta.Error != nil {
		return respuesta.Error
	}

	return nil
}
