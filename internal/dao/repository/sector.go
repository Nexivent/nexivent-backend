package repository

import (
	"time"

	"github.com/Loui27/nexivent-backend/internal/dao/model"
	"github.com/Loui27/nexivent-backend/logging"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *Sector) ModificarSectorPorCampos(
	id int64,
	sectorTipo *string,
	totalEntradas *int,
	cantVendidas *int, // úsalo solo para correcciones administrativas
	estado *int16,
	usuarioModificacion *int64,
	fechaModificacion *time.Time,
) (*model.Sector, error) {

	if id <= 0 {
		return nil, gorm.ErrInvalidData
	}

	updates := map[string]any{}
	if sectorTipo != nil {
		updates["sector_tipo"] = *sectorTipo
	}
	if totalEntradas != nil {
		updates["total_entradas"] = *totalEntradas
	}
	if cantVendidas != nil {
		updates["cant_vendidas"] = *cantVendidas
	}
	if estado != nil {
		updates["estado"] = *estado
	}
	if usuarioModificacion != nil {
		updates["usuario_modificacion"] = *usuarioModificacion
	}
	if fechaModificacion != nil {
		updates["fecha_modificacion"] = *fechaModificacion
	}

	var s model.Sector
	if len(updates) == 0 {
		if err := r.PostgresqlDB.First(&s, "sector_id = ?", id).Error; err != nil {
			r.logger.Errorf("ModificarSectorCampos (sin cambios) id=%d: %v", id, err)
			return nil, err
		}
		return &s, nil
	}

	res := r.PostgresqlDB.
		Model(&s).
		Clauses(clause.Returning{}).
		Where("sector_id = ?", id).
		Updates(updates)

	if res.Error != nil {
		r.logger.Errorf("ModificarSectorCampos id=%d: %v", id, res.Error)
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &s, nil
}
