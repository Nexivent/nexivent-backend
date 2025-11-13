package repository

import (
	"time"

	"github.com/Loui27/nexivent-backend/internal/dao/model"
	"github.com/Loui27/nexivent-backend/logging"
	"gorm.io/gorm"
)

type Comentario struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewComentariosController(
	logger logging.Logger,
	postgresqlDB *gorm.DB,
) *Comentario {
	return &Comentario{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// CrearComentario: inserta un comentario (FKs deben existir en BD).
func (r *Comentario) CrearComentario(c *model.Comentario) error {
	if c == nil {
		return gorm.ErrInvalidData
	}
	if err := r.PostgresqlDB.Create(c).Error; err != nil {
		r.logger.Errorf("CrearComentario: %v", err)
		return err
	}
	return nil
}

// ObtenerComentarioPorID: retorna un comentario por su ID.
func (r *Comentario) ObtenerComentarioPorID(id int64) (*model.Comentario, error) {
	var c model.Comentario
	if err := r.PostgresqlDB.First(&c, "comentario_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// ListarComentariosActivosPorEvento: comentarios estado=1, ordenados por fecha (ASC).
// func (r *Comentario) ListarComentariosActivosPorEvento(
// 	eventoID int64,
// 	limit int,
// 	offset int,
// ) ([]model.Comentario, error) {

// 	var list []model.Comentario
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
// 		r.logger.Errorf("ListarComentariosActivosPorEvento(evento=%d): %v", eventoID, err)
// 		return nil, err
// 	}
// 	return list, nil
// }

// ContarComentariosActivosPorEvento: total de comentarios activos del evento.
func (r *Comentario) ContarComentariosActivosPorEvento(eventoID int64) (int64, error) {
	var count int64
	res := r.PostgresqlDB.
		Table("comentario").
		Where("evento_id = ? AND estado = 1", eventoID).
		Count(&count)
	if res.Error != nil {
		r.logger.Errorf("ContarComentariosActivosPorEvento(evento=%d): %v", eventoID, res.Error)
		return 0, res.Error
	}
	return count, nil
}

// DesactivarComentario: estado = 0 (soft off).
func (r *Comentario) DesactivarComentario(
	id int64,
	usuarioModificacion *int64,
	fechaModificacion *time.Time,
) error {

	if id <= 0 {
		return gorm.ErrInvalidData
	}

	updates := map[string]any{
		"estado": int16(0),
	}
	if usuarioModificacion != nil {
		updates["usuario_modificacion"] = *usuarioModificacion
	}
	if fechaModificacion != nil {
		updates["fecha_modificacion"] = *fechaModificacion
	}

	res := r.PostgresqlDB.
		Table("comentario").
		Where("comentario_id = ?", id).
		Updates(updates)

	if res.Error != nil {
		r.logger.Errorf("DesactivarComentario(id=%d): %v", id, res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Esto si los eventos se reactivan para otra fecha y la gente pueda ver comentarios de antiguos eventos??
// ActivarComentario: estado = 1 (soft on) — por si necesitas reactivar.
func (r *Comentario) ActivarComentario(
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
		Table("comentario").
		Where("comentario_id = ?", id).
		Updates(updates)

	if res.Error != nil {
		r.logger.Errorf("ActivarComentario(id=%d): %v", id, res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// VerificarEventoExiste: true si el evento existe (sin validar flags/estado de workflow).
func (r *Comentario) VerificarEventoExiste(eventoID int64) (bool, error) {
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
func (r *Comentario) VerificarUsuarioExiste(usuarioID int64) (bool, error) {
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

func (c *Comentario) ListarComentariosPorIdEvento(eventoID int64) ([]*model.Comentario, error) {
	var comentarios []*model.Comentario
	respuesta := c.PostgresqlDB.Where("evento_id = ?", eventoID).
		Find(&comentarios)
	if respuesta.Error != nil {
		return nil, respuesta.Error
	}
	return comentarios, nil
}
