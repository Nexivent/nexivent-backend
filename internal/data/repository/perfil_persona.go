package repository

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PerfilDePersona struct {
	DB *gorm.DB
}

func NewPerfilDePersonaController(
	postgresqlDB *gorm.DB, // (opcional: corregí el nombre)
) *PerfilDePersona {
	return &PerfilDePersona{
		DB: postgresqlDB,
	}
}

func (p *PerfilDePersona) CrearPerfilDePersona(perfil *model.PerfilDePersona) error {
	if err := p.DB.Create(perfil).Error; err != nil {
		return err
	}
	return nil
}

func (p *PerfilDePersona) ActualizarPerfilDePersona(perfil *model.PerfilDePersona) error {
	if err := p.DB.Save(perfil).Error; err != nil {
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
		if err := r.DB.First(&p, "perfil_de_persona_id = ?", id).Error; err != nil {
			return nil, err
		}
		return &p, nil
	}

	res := r.DB.
		Model(&p).
		Clauses(clause.Returning{}).
		Where("perfil_de_persona_id = ?", id).
		Updates(updates)

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &p, nil
}

func (pd *PerfilDePersona) ListarPerfilPersonaPorEventoID(eventoID uint64) ([]model.PerfilDePersona, error) {
	var list []model.PerfilDePersona
	if err := pd.DB.
		Where("evento_id = ?", eventoID).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
