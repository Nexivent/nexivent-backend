package repository

import (
	"time"

	"github.com/Loui27/nexivent-backend/internal/dao/model"
	"github.com/Loui27/nexivent-backend/logging"
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

	// 0) Si ya está ACTIVO, devolver tal cual (previene duplicados activos)
	{
		var existente model.RolUsuario
		check := r.PostgresqlDB.
			Where("usuario_id = ? AND rol_id = ? AND estado = 1", usuarioID, rolID).
			First(&existente)
		if check.Error == nil {
			return &existente, nil
		}
		if check.Error != nil && check.Error != gorm.ErrRecordNotFound {
			return nil, check.Error
		}
	}

	// 1) Si existe INACTIVO (estado=0) -> reactivar
	result := r.PostgresqlDB.
		Model(&model.RolUsuario{}).
		Where("usuario_id = ? AND rol_id = ? AND estado = 0", usuarioID, rolID).
		Updates(map[string]any{
			"estado":               int16(1),
			"usuario_modificacion": createdBy, // campo puntero en modelo, map con valor OK
			"fecha_modificacion":   now,       // idem
		})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		var ru model.RolUsuario
		get := r.PostgresqlDB.Where("usuario_id = ? AND rol_id = ?", usuarioID, rolID).First(&ru)
		if get.Error != nil {
			return nil, get.Error
		}
		return &ru, nil
	}

	// 2) No existía -> insertar NUEVO (tipos coherentes con el modelo)
	ru := &model.RolUsuario{
		RolID:               rolID,
		UsuarioID:           usuarioID,
		UsuarioCreacion:     createdBy,  // int64 (no nulo)
		FechaCreacion:       now,        // time.Time (no nulo)
		UsuarioModificacion: &createdBy, // *int64
		FechaModificacion:   &now,       // *time.Time
		Estado:              1,
	}
	insert := r.PostgresqlDB.Create(ru)
	if insert.Error != nil {
		r.logger.Errorf("AsignarRolAUsuario uid=%v rol=%v: %v", usuarioID, rolID, insert.Error)
		return nil, insert.Error
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
			"estado":               int16(0),
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

// Listar roles activos de un usuario (con preload del Rol)
func (r *RolUsuarioRepo) ListarRolesDeUsuario(usuarioID int64) ([]*model.RolUsuario, error) {
	asignaciones := []*model.RolUsuario{}
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
