package repository

import (
	"errors"
	"time"

	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Rol struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewRolesController(logger logging.Logger, postgresqlDB *gorm.DB) *Rol {
	return &Rol{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// Listar todos los roles
func (r *Rol) ObtenerRoles() ([]*model.Rol, error) {
	roles := []*model.Rol{}
	result := r.PostgresqlDB.Find(&roles)
	if result.Error != nil {
		r.logger.Errorf("ObtenerRoles: %v", result.Error)
		return nil, result.Error
	}
	return roles, nil
}

// Crear rol (sugiere UNIQUE(nombre) en BD)
func (r *Rol) CrearRol(rol *model.Rol) error {
	if rol == nil {
		return errors.New("CrearRol: rol es nil")
	}
	result := r.PostgresqlDB.Create(rol)
	if result.Error != nil {
		r.logger.Errorf("CrearRol: %v", result.Error)
		return result.Error
	}
	return nil
}

// Actualizar nombre con auditoría
func (r *Rol) ActualizarRol(
	id int64,
	nombre *string,
	updatedBy int64,
) (*model.Rol, error) {
	updateFields := map[string]any{
		"usuario_modificacion": updatedBy,
		"fecha_modificacion":   time.Now(),
	}
	if nombre != nil {
		updateFields["nombre"] = *nombre
	}

	var rol model.Rol
	// Si solo llegó auditoría, devolvemos el registro actual
	if len(updateFields) == 2 {
		if err := r.PostgresqlDB.First(&rol, "rol_id = ?", id).Error; err != nil {
			r.logger.Errorf("ActualizarRol(sin cambios) id=%v: %v", id, err)
			return nil, err
		}
		return &rol, nil
	}

	result := r.PostgresqlDB.Model(&rol).
		Clauses(clause.Returning{}).
		Where("rol_id = ?", id).
		Updates(updateFields)

	if result.Error != nil {
		r.logger.Errorf("ActualizarRol id=%v: %v", id, result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &rol, nil
}

// Utilidad: obtener rol por nombre (para validar duplicados)
func (r *Rol) ObtenerRolPorNombre(nombre string) (*model.Rol, error) {
	var rol model.Rol
	result := r.PostgresqlDB.
		Where("nombre = ?", nombre).
		First(&rol)
	if result.Error != nil {
		return nil, result.Error
	}
	return &rol, nil
}

func (r *Rol) ObtenerRolPorID(idRol int64) (*model.Rol, error) {
	var rol model.Rol
	result := r.PostgresqlDB.
		Where("rol_id = ?", idRol).
		First(&rol)
	if result.Error != nil {
		return nil, result.Error
	}
	return &rol, nil
}

func (r *Rol) ObtenerRolesDeUsuario(usuarioID int64) ([]*model.Rol, error) {
	r.logger.Infof("📊 [REPO] Buscando roles para usuario ID: %d", usuarioID)

    var roles []*model.Rol

    res := r.PostgresqlDB.
        Table("rol").
        Select("rol.rol_id as rol_id, rol.nombre as nombre, rol.usuario_creacion, rol.fecha_creacion, rol.usuario_modificacion, rol.fecha_modificacion").
        Joins("INNER JOIN rol_usuario ru ON ru.rol_id = rol.rol_id AND ru.estado = 1").
        Where("ru.usuario_id = ?", usuarioID).
        Scan(&roles)

    if res.Error != nil {
        r.logger.Errorf("❌ [REPO] Error obteniendo roles: %v", res.Error)
        return nil, res.Error
    }

    r.logger.Infof("✅ [REPO] Roles encontrados: %d", len(roles))
    for i, rol := range roles {
        r.logger.Infof("   [REPO] Rol %d: ID=%d, Nombre=%s", i+1, rol.ID, rol.Nombre)
    }

    return roles, nil
}
