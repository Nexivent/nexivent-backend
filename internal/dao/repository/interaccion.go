package repository

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	util "github.com/Nexivent/nexivent-backend/internal/dao/model/util"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
)

type Interaccion struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewInteraccionController(
	logger logging.Logger,
	postgresqlDB *gorm.DB,
) *Interaccion {
	return &Interaccion{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// CrearInteraccion: inserta un Interaccion (FKs deben existir en BD).
func (r *Interaccion) CrearInteraccion(c *model.Interaccion) error {
	if c == nil {
		return gorm.ErrInvalidData
	}
	if err := r.PostgresqlDB.Create(c).Error; err != nil {
		r.logger.Errorf("CrearInteraccion: %v", err)
		return err
	}
	return nil
}

// ObtenerInteraccionPorID: retorna un Interaccion por su ID.
func (r *Interaccion) ObtenerInteraccionPorID(id int64) (*model.Interaccion, error) {
	var c model.Interaccion
	if err := r.PostgresqlDB.First(&c, "interaccion_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// ListarInteraccionsActivosPorEvento: Interaccions estado=1, ordenados por fecha (ASC).
// func (r *Interaccion) ListarInteraccionsActivosPorEvento(
// 	eventoID int64,
// 	limit int,
// 	offset int,
// ) ([]model.Interaccion, error) {

// 	var list []model.Interaccion
// 	q := r.PostgresqlDB.
// 		Where("evento_id = ? AND estado = 1", eventoID).
// 		Order("fecha_creacion ASC")

// 	if limit > 0 {
// 		q = q.Limit(limit)
// 	}
// 	if offset > 0 {
// 		q = q.Offset(offset)
// 	}

// 	if err := q.Find(&list).Error; err != nil {
// 		r.logger.Errorf("ListarInteraccionsActivosPorEvento(evento=%d): %v", eventoID, err)
// 		return nil, err
// 	}
// 	return list, nil
// }

// ContarInteraccionsActivosPorEvento: total de Interaccions activos del evento.
func (r *Interaccion) ContarInteraccionsActivosPorEvento(eventoID int64) (int64, error) {
	var count int64
	res := r.PostgresqlDB.
		Table("interaccion").
		Where("evento_id = ? AND estado = 1", eventoID).
		Count(&count)
	if res.Error != nil {
		r.logger.Errorf("ContarInteraccionsActivosPorEvento(evento=%d): %v", eventoID, res.Error)
		return 0, res.Error
	}
	return count, nil
}

// DesactivarInteraccion: estado = 0 (soft off).
func (r *Interaccion) DesactivarInteraccion(
	id int64,
	usuarioModificacion *int64,
	fechaModificacion *time.Time,
) error {

	if id <= 0 {
		return gorm.ErrInvalidData
	}

	updates := map[string]any{
		"estado": util.Inactivo,
	}
	if usuarioModificacion != nil {
		updates["usuario_modificacion"] = *usuarioModificacion
	}
	if fechaModificacion != nil {
		updates["fecha_modificacion"] = *fechaModificacion
	}

	res := r.PostgresqlDB.
		Table("interaccion").
		Where("interaccion_id = ?", id).
		Updates(updates)

	if res.Error != nil {
		r.logger.Errorf("DesactivarInteraccion(id=%d): %v", id, res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Esto si los eventos se reactivan para otra fecha y la gente pueda ver Interaccions de antiguos eventos??
// ActivarInteraccion: estado = 1 (soft on) — por si necesitas reactivar.
func (r *Interaccion) ActivarInteraccion(
	id int64,
	usuarioModificacion *int64,
	fechaModificacion *time.Time,
) error {

	if id <= 0 {
		return gorm.ErrInvalidData
	}

	updates := map[string]any{
		"estado": int16(1),
	}
	if usuarioModificacion != nil {
		updates["usuario_modificacion"] = *usuarioModificacion
	}
	if fechaModificacion != nil {
		updates["fecha_modificacion"] = *fechaModificacion
	}

	res := r.PostgresqlDB.
		Table("interaccion").
		Where("interaccion_id = ?", id).
		Updates(updates)

	if res.Error != nil {
		r.logger.Errorf("ActivarInteraccion(id=%d): %v", id, res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// VerificarEventoExiste: true si el evento existe (sin validar flags/estado de workflow).
func (r *Interaccion) VerificarEventoExiste(eventoID int64) (bool, error) {
	var count int64
	res := r.PostgresqlDB.
		Table("evento").
		Where("evento_id = ?", eventoID).
		Count(&count)
	if res.Error != nil {
		r.logger.Errorf("VerificarEventoExiste(evento=%d): %v", eventoID, res.Error)
		return false, res.Error
	}
	return count == 1, nil
}

// VerificarUsuarioExiste: true si el usuario existe (sin validar estado).
func (r *Interaccion) VerificarUsuarioExiste(usuarioID int64) (bool, error) {
	var count int64
	res := r.PostgresqlDB.
		Table("usuario").
		Where("usuario_id = ?", usuarioID).
		Count(&count)
	if res.Error != nil {
		r.logger.Errorf("VerificarUsuarioExiste(usuario=%d): %v", usuarioID, res.Error)
		return false, res.Error
	}
	return count == 1, nil
}

func (c *Interaccion) ListarInteraccionsPorIdEvento(eventoID int64) ([]*model.Interaccion, error) {
	var Interaccions []*model.Interaccion
	respuesta := c.PostgresqlDB.Where("evento_id = ?", eventoID).
		Find(&Interaccions)
	if respuesta.Error != nil {
		return nil, respuesta.Error
	}
	return Interaccions, nil
}

func (c *Interaccion) ObtenerInteraccionesEventoUsuario(eventoId int64, usuarioId int64) (*model.Interaccion, error) {
	var interaccion *model.Interaccion

	respuesta := c.PostgresqlDB.
		Table("interaccion").
		Where("evento_id = ? AND usuario_id = ?", eventoId, usuarioId).
		Find(&interaccion)
	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	return interaccion, nil
}

func (c *Interaccion) ActualizarInteracciones(interacciones model.Interaccion) error {
	respuesta := c.PostgresqlDB.Table("interaccion").Where("interaccion_id = ?", interacciones.ID).Update("tipo", interacciones.Tipo)
	if respuesta != nil {
		return respuesta.Error
	}
	return nil
}
