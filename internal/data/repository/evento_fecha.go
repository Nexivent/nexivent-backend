package repository

import (
	"database/sql"

	"github.com/Nexivent/nexivent-backend/internal/data/model"
	"gorm.io/gorm"
)

type EventoFecha struct {
	DB *gorm.DB
}

func (e *EventoFecha) CrearEventoFecha(EventoFecha *model.EventoFecha) error {
	respuesta := e.DB.Create(EventoFecha)

	if respuesta.Error != nil {
		return respuesta.Error
	}

	return nil
}

func (e *EventoFecha) ListarEventoFechasPorEventoID(eventoID uint64) ([]model.EventoFecha, error) {
	var list []model.EventoFecha
	if err := e.DB.
		Where("evento_id = ?", eventoID).
		Order("hora_inicio ASC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// ListarEventoFechasActivasPorEventoID: solo estado = 1
func (e *EventoFecha) ListarEventoFechasActivasPorEventoID(eventoID uint64) ([]model.EventoFecha, error) {
	var list []model.EventoFecha
	if err := e.DB.
		Where("evento_id = ? AND estado = 1", eventoID).
		Order("hora_inicio ASC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// (Opcional) Con preload de la tabla fecha (por si necesitas el DATE ya cargado)
func (e *EventoFecha) ListarEventoFechasActivasConFecha(eventoID uint64) ([]model.EventoFecha, error) {
	var list []model.EventoFecha
	if err := e.DB.
		Preload("Fecha"). // requiere que model.EventoFecha tenga: Fecha *model.Fecha `gorm:"foreignKey:FechaID;references:fecha_id"`
		Where("evento_id = ? AND estado = 1", eventoID).
		Order("hora_inicio ASC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// VerificarEventoExiste: true si existe evento.
func (e *EventoFecha) VerificarEventoExiste(eventoID uint64) (bool, error) {
	var count int64
	res := e.DB.Table("evento").Where("evento_id = ?", eventoID).Count(&count)
	if res.Error != nil {
		return false, res.Error
	}
	return count == 1, nil
}

// VerificarFechaExiste: true si existe fila en tabla fecha.
func (e *EventoFecha) VerificarFechaExiste(fechaID uint64) (bool, error) {
	var count int64
	res := e.DB.Table("fecha").Where("fecha_id = ?", fechaID).Count(&count)
	if res.Error != nil {
		return false, res.Error
	}
	return count == 1, nil
}

// VerificarEventoFechaDuplicada: true si ya existe (evento_id, fecha_id, hora_inicio).
func (e *EventoFecha) VerificarEventoFechaDuplicada(eventoID, fechaID uint64, horaInicio any) (bool, error) {
	// horaInicio puede ser time.Time o string con formato compatible con PG (el BO pasa el valor exacto)
	var count int64
	res := e.DB.
		Table("evento_fecha").
		Where("evento_id = ? AND fecha_id = ? AND hora_inicio = ?", eventoID, fechaID, horaInicio).
		Count(&count)
	if res.Error != nil {
		return false, res.Error
	}
	return count > 0, nil
}

// VerificarEventoFechaActiva: true si estado = 1.
func (e *EventoFecha) VerificarEventoFechaActiva(eventoFechaID uint64) (bool, error) {
	var estado int16
	err := e.DB.
		Table("evento_fecha").
		Select("estado").
		Where("evento_fecha_id = ?", eventoFechaID).
		Row().
		Scan(&estado)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, gorm.ErrRecordNotFound
		}
		return false, err
	}
	return estado == 1, nil
}

// VerificarEventoFechaPerteneceAEvento: true si el evento_fecha pertenece al evento.
func (e *EventoFecha) VerificarEventoFechaPerteneceAEvento(eventoFechaID, eventoID uint64) (bool, error) {
	var count int64
	res := e.DB.
		Table("evento_fecha").
		Where("evento_fecha_id = ? AND evento_id = ?", eventoFechaID, eventoID).
		Count(&count)
	if res.Error != nil {
		return false, res.Error
	}
	return count == 1, nil
}

// ContarEventoFechasActivas: total de fechas activas de un evento (útil para UI/validaciones).
func (e *EventoFecha) ContarEventoFechasActivas(eventoID uint64) (int64, error) {
	var count int64
	res := e.DB.
		Table("evento_fecha").
		Where("evento_id = ? AND estado = 1", eventoID).
		Count(&count)
	if res.Error != nil {
		return 0, res.Error
	}
	return count, nil
}

// ActivarEventoFecha cambia el estado de un evento_fecha a activo (1)
func (e *EventoFecha) ActivarEventoFecha(eventoFechaID uint64) error {
	res := e.DB.
		Table("evento_fecha").
		Where("evento_fecha_id = ?", eventoFechaID).
		Update("estado", int16(1))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DesactivarEventoFecha cambia el estado de un evento_fecha a inactivo (0)
func (e *EventoFecha) DesactivarEventoFecha(eventoFechaID uint64) error {
	res := e.DB.
		Table("evento_fecha").
		Where("evento_fecha_id = ?", eventoFechaID).
		Update("estado", int16(0))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}