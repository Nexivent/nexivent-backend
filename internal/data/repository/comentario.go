package repository

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/model"
	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"gorm.io/gorm"
)

type Comentario struct {
	DB *gorm.DB
}

// CrearComentario: inserta un comentario (FKs deben existir en BD).
func (c *Comentario) CrearComentario(comentario *model.Comentario) error {
	if comentario == nil {
		return gorm.ErrInvalidData
	}
	if err := c.DB.Create(comentario).Error; err != nil {
		return err
	}
	return nil
}

// ObtenerComentarioPorID: retorna un comentario por su ID.
func (c *Comentario) ObtenerComentarioPorID(id int64) (*model.Comentario, error) {
	var comentario model.Comentario
	if err := c.DB.First(&comentario, "comentario_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &comentario, nil
}

// ListarComentariosActivosPorEvento: comentarios estado=1, ordenados por fecha (ASC).
// func (c *Comentario) ListarComentariosActivosPorEvento(
// 	eventoID int64,
// 	limit int,
// 	offset int,
// ) ([]model.Comentario, error) {

// 	var list []model.Comentario
// 	q := r.DB.
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
func (c *Comentario) ContarComentariosActivosPorEvento(eventoID int64) (int64, error) {
	var count int64
	res := c.DB.
		Table("comentario").
		Where("evento_id = ? AND estado = 1", eventoID).
		Count(&count)
	if res.Error != nil {
		return 0, res.Error
	}
	return count, nil
}

// DesactivarComentario: estado = 0 (soft off).
func (c *Comentario) DesactivarComentario(
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

	res := c.DB.
		Table("comentario").
		Where("comentario_id = ?", id).
		Updates(updates)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Esto si los eventos se reactivan para otra fecha y la gente pueda ver comentarios de antiguos eventos??
// ActivarComentario: estado = 1 (soft on) — por si necesitas reactivar.
func (c *Comentario) ActivarComentario(
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

	res := c.DB.
		Table("comentario").
		Where("comentario_id = ?", id).
		Updates(updates)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// VerificarEventoExiste: true si el evento existe (sin validar flags/estado de workflow).
func (c *Comentario) VerificarEventoExiste(eventoID int64) (bool, error) {
	var count int64
	res := c.DB.
		Table("evento").
		Where("evento_id = ?", eventoID).
		Count(&count)
	if res.Error != nil {
		return false, res.Error
	}
	return count == 1, nil
}

// VerificarUsuarioExiste: true si el usuario existe (sin validar estado).
func (c *Comentario) VerificarUsuarioExiste(usuarioID int64) (bool, error) {
	var count int64
	res := c.DB.
		Table("usuario").
		Where("usuario_id = ?", usuarioID).
		Count(&count)
	if res.Error != nil {
		return false, res.Error
	}
	return count == 1, nil
}

func (c *Comentario) ListarComentariosPorIdEvento(eventoID uint64) ([]model.Comentario, error) {
	var comentarios []model.Comentario
	respuesta := c.DB.Where("evento_id = ?", eventoID).
		Find(&comentarios)
	if respuesta.Error != nil {
		return nil, respuesta.Error
	}
	return comentarios, nil
}
