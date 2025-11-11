package repository

import (
	"time"

	"github.com/Loui27/nexivent-backend/internal/dao/model"
	"github.com/Loui27/nexivent-backend/logging"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PerfilDePersona struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewPerfilDePersonaController(
	logger logging.Logger,
	postgresqlDB *gorm.DB, // (opcional: corregí el nombre)
) *PerfilDePersona {
	return &PerfilDePersona{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

func (p *PerfilDePersona) CrearPerfilDePersona(perfil *model.PerfilDePersona) error {
	if err := p.PostgresqlDB.Create(perfil).Error; err != nil {
		return err
	}
	return nil
}

func (p *PerfilDePersona) ActualizarPerfilDePersona(perfil *model.PerfilDePersona) error {
	if err := p.PostgresqlDB.Save(perfil).Error; err != nil {
		return err
	}
	return nil
}

func (r *PerfilDePersona) ModificarPerfilDePersonaPorCampos(
	id int64,
	eventoID *int64,
	nombre *string,
	estado *int16,
	usuarioModificacion *int64,
	fechaModificacion *time.Time,
) (*model.PerfilDePersona, error) {

	if id <= 0 {
		return nil, gorm.ErrInvalidData
	}

	updates := map[string]any{}
	if eventoID != nil {
		updates["evento_id"] = *eventoID
	}
	if nombre != nil {
		updates["nombre"] = *nombre
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

	var p model.PerfilDePersona
	if len(updates) == 0 {
		if err := r.PostgresqlDB.First(&p, "perfil_de_persona_id = ?", id).Error; err != nil {
			r.logger.Errorf("ModificarPerfilDePersonaCampos (sin cambios) id=%d: %v", id, err)
			return nil, err
		}
		return &p, nil
	}

	res := r.PostgresqlDB.
		Model(&p).
		Clauses(clause.Returning{}).
		Where("perfil_de_persona_id = ?", id).
		Updates(updates)

	if res.Error != nil {
		r.logger.Errorf("ModificarPerfilDePersonaCampos id=%d: %v", id, res.Error)
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &p, nil
}
