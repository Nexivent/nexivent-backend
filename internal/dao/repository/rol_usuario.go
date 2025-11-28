package repository

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	util "github.com/Nexivent/nexivent-backend/internal/dao/model/util"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RolUsuarioRepo struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewRolUsuarioController(logger logging.Logger, postgresqlDB *gorm.DB) *RolUsuarioRepo {
	return &RolUsuarioRepo{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// Asignar rol a usuario (si ya existe y está activo, no duplica; si existe inactivo, puedes reactivarlo)
func (r *RolUsuarioRepo) AsignarRolAUsuario(
	usuarioID int64,
	rolID int64,
	createdBy int64,
) (*model.RolUsuario, error) {

	now := time.Now()
	ru := &model.RolUsuario{
		RolID:               rolID,
		UsuarioID:           usuarioID,
		UsuarioCreacion:     &createdBy,
		FechaCreacion:       now,
		UsuarioModificacion: &createdBy,
		FechaModificacion:   &now,
		Estado:              1,
	}

	// INSERT ... ON CONFLICT (usuario_id, rol_id) DO UPDATE SET estado=1, audit...
	err := r.PostgresqlDB.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "usuario_id"},
				{Name: "rol_id"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"estado":               int16(1),
				"usuario_modificacion": createdBy,
				"fecha_modificacion":   now,
			}),
		}).
		Create(ru).Error

	if err != nil {
		return nil, err
	}
	return ru, nil
}

// Desactivar asignación (pasar 1 a 0) y no hacer nada si ya estaba en 0
func (r *RolUsuarioRepo) QuitarRolDeUsuario(
	usuarioID int64,
	rolID int64,
	updatedBy int64,
) error {
	result := r.PostgresqlDB.
		Model(&model.RolUsuario{}).
		Where("usuario_id = ? AND rol_id = ? AND estado = 1", usuarioID, rolID).
		Updates(map[string]any{
			"estado":               util.Inactivo,
			"usuario_modificacion": updatedBy,
			"fecha_modificacion":   time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// borrado fisico
func (r *RolUsuarioRepo) BorrarRolDeUsuario(
	usuarioID int64,
	rolID int64,
	updatedBy int64,
) error {
	result := r.PostgresqlDB.
		Where("usuario_id = ? AND rol_id = ?", usuarioID, rolID).
		Delete(&model.RolUsuario{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Listar roles activos de un usuario (con preload del Rol)
func (r *RolUsuarioRepo) ListarRolesDeUsuario(usuarioID int64) ([]model.RolUsuario, error) {
	asignaciones := []model.RolUsuario{}
	result := r.PostgresqlDB.
		Preload("Rol").
		Where("usuario_id = ? AND estado = 1", usuarioID).
		Find(&asignaciones)
	if result.Error != nil {
		return nil, result.Error
	}
	return asignaciones, nil
}

// Existe asignación activa
func (r *RolUsuarioRepo) ExisteRolUsuario(usuarioID int64, rolID int64) (bool, error) {
	var count int64
	result := r.PostgresqlDB.
		Model(&model.RolUsuario{}).
		Where("usuario_id = ? AND rol_id = ? AND estado = 1", usuarioID, rolID).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// Obtener usuarios por un rol específico
func (ru *RolUsuarioRepo) ObtenerUsuariosPorRol(rolId int64) ([]model.Usuario, error) {
	var usuarios []model.Usuario
	respuesta := ru.PostgresqlDB.
		Table("usuario").
		Joins("JOIN rol_usuario r ON r.usuario_id = usuario.usuario_id").
		Where("r.rol_id = ? AND r.estado = 1 ", rolId).
		Find(&usuarios)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	return usuarios, nil
}

// Obtener TODOS los usuarios activos (sin importar sus roles)
func (ru *RolUsuarioRepo) ObtenerTodosLosUsuariosActivos() ([]model.Usuario, error) {
	var usuarios []model.Usuario
	respuesta := ru.PostgresqlDB.
		//Where("estado = 1").
		Find(&usuarios)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	return usuarios, nil
}

// Actualizar estado de una asignación por id (útil para admin)
func (r *RolUsuarioRepo) ActualizarRolUsuarioEstado(
	rolUsuarioID int64,
	estado int16,
	updatedBy int64,
) (*model.RolUsuario, error) {
	var ru model.RolUsuario
	result := r.PostgresqlDB.Model(&ru).
		Clauses(clause.Returning{}).
		Where("rol_usuario_id = ?", rolUsuarioID).
		Updates(map[string]any{
			"estado":               estado,
			"usuario_modificacion": updatedBy,
			"fecha_modificacion":   time.Now(),
		})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &ru, nil
}

// Helper interno
func ptrTime(t time.Time) *time.Time { return &t }
