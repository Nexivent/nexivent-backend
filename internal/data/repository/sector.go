package repository

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Sector struct {
	DB *gorm.DB
}

func (s *Sector) CrearSector(sector *model.Sector) error {
	resultado := s.DB.Create(sector)

	if resultado.Error != nil {
		return resultado.Error
	}

	return nil
}

func (s *Sector) ActualizarSector(sector *model.Sector) error {
	respuesta := s.DB.Save(sector)

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
		if err := r.DB.First(&s, "sector_id = ?", id).Error; err != nil {
			return nil, err
		}
		return &s, nil
	}

	res := r.DB.
		Model(&s).
		Clauses(clause.Returning{}).
		Where("sector_id = ?", id).
		Updates(updates)

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &s, nil
}
