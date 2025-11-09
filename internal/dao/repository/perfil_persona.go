package repository

import (
	"github.com/Loui27/nexivent-backend/internal/dao/model"
	"github.com/Loui27/nexivent-backend/logging"
	"gorm.io/gorm"
)

type PerfilDePersona struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewPerfilDePersonaController(
	logger logging.Logger,
	postgreesqlDB *gorm.DB,
) *PerfilDePersona {
	return &PerfilDePersona{
		logger:       logger,
		PostgresqlDB: postgreesqlDB,
	}
}

func (p *PerfilDePersona) CrearSector(PerfilDePersona *model.PerfilDePersona) error {
	resultado := p.PostgresqlDB.Create(PerfilDePersona)

	if resultado.Error != nil {
		return resultado.Error
	}

	return nil
}

func (p *PerfilDePersona) ActualizarSector(PerfilDePersona *model.PerfilDePersona) error {
	respuesta := p.PostgresqlDB.Save(PerfilDePersona)

	if respuesta.Error != nil {
		return respuesta.Error
	}

	return nil
}
