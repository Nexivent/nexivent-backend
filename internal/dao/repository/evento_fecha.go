package repository

import (
	"database/sql"

	"github.com/Loui27/nexivent-backend/internal/dao/model"
	"github.com/Loui27/nexivent-backend/logging"
	"gorm.io/gorm"
	// "github.com/Loui27/nexivent-backend/internal/dao/model"
)

type EventoFecha struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewEventoFechaController(
	logger logging.Logger,
	postgreesqlDB *gorm.DB,
) *EventoFecha {
	return &EventoFecha{
		logger:       logger,
		PostgresqlDB: postgreesqlDB,
	}
}

func (e *EventoFecha) CrearEventoFecha(EventoFecha *model.EventoFecha) error {
	respuesta := e.PostgresqlDB.Create(EventoFecha)

	if respuesta.Error != nil {
		return respuesta.Error
	}

	return nil
}

// ListarEventoFechasPorEventoID: devuelve TODAS las filas por evento_id (sin filtrar estado)
func (r *EventoFecha) ListarEventoFechasPorEventoID(eventoID int64) ([]model.EventoFecha, error) {
	var list []model.EventoFecha
	if err := r.PostgresqlDB.
		Where("evento_id = ?", eventoID).
		Order("hora_inicio ASC").
		Find(&list).Error; err != nil {
		r.logger.Errorf("ListarEventoFechasPorEventoID(%d): %v", eventoID, err)
		return nil, err
	}
	return list, nil
}

// ListarEventoFechasActivasPorEventoID: solo estado = 1
func (r *EventoFecha) ListarEventoFechasActivasPorEventoID(eventoID int64) ([]model.EventoFecha, error) {
	var list []model.EventoFecha
	if err := r.PostgresqlDB.
		Where("evento_id = ? AND estado = 1", eventoID).
		Order("hora_inicio ASC").
		Find(&list).Error; err != nil {
		r.logger.Errorf("ListarEventoFechasActivasPorEventoID(%d): %v", eventoID, err)
		return nil, err
	}
	return list, nil
}

// (Opcional) Con preload de la tabla fecha (por si necesitas el DATE ya cargado)
func (r *EventoFecha) ListarEventoFechasActivasConFecha(eventoID int64) ([]model.EventoFecha, error) {
	var list []model.EventoFecha
	if err := r.PostgresqlDB.
		Preload("Fecha"). // requiere que model.EventoFecha tenga: Fecha *model.Fecha `gorm:"foreignKey:FechaID;references:fecha_id"`
		Where("evento_id = ? AND estado = 1", eventoID).
		Order("hora_inicio ASC").
		Find(&list).Error; err != nil {
		r.logger.Errorf("ListarEventoFechasActivasConFecha(%d): %v", eventoID, err)
		return nil, err
	}
	return list, nil
}

// VerificarEventoExiste: true si existe evento.
func (r *EventoFecha) VerificarEventoExiste(eventoID int64) (bool, error) {
	var count int64
	res := r.PostgresqlDB.Table("evento").Where("evento_id = ?", eventoID).Count(&count)
	if res.Error != nil {
		r.logger.Errorf("VerificarEventoExiste(%d): %v", eventoID, res.Error)
		return false, res.Error
	}
	return count == 1, nil
}

// VerificarFechaExiste: true si existe fila en tabla fecha.
func (r *EventoFecha) VerificarFechaExiste(fechaID int64) (bool, error) {
	var count int64
	res := r.PostgresqlDB.Table("fecha").Where("fecha_id = ?", fechaID).Count(&count)
	if res.Error != nil {
		r.logger.Errorf("VerificarFechaExiste(%d): %v", fechaID, res.Error)
		return false, res.Error
	}
	return count == 1, nil
}

// VerificarEventoFechaDuplicada: true si ya existe (evento_id, fecha_id, hora_inicio).
func (r *EventoFecha) VerificarEventoFechaDuplicada(eventoID, fechaID int64, horaInicio any) (bool, error) {
	// horaInicio puede ser time.Time o string con formato compatible con PG (el BO pasa el valor exacto)
	var count int64
	res := r.PostgresqlDB.
		Table("evento_fecha").
		Where("evento_id = ? AND fecha_id = ? AND hora_inicio = ?", eventoID, fechaID, horaInicio).
		Count(&count)
	if res.Error != nil {
		r.logger.Errorf("VerificarEventoFechaDuplicada(e=%d,f=%d): %v", eventoID, fechaID, res.Error)
		return false, res.Error
	}
	return count > 0, nil
}

// VerificarEventoFechaActiva: true si estado = 1.
func (r *EventoFecha) VerificarEventoFechaActiva(eventoFechaID int64) (bool, error) {
	var estado int16
	err := r.PostgresqlDB.
		Table("evento_fecha").
		Select("estado").
		Where("evento_fecha_id = ?", eventoFechaID).
		Row().
		Scan(&estado)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, gorm.ErrRecordNotFound
		}
		r.logger.Errorf("VerificarEventoFechaActiva(%d): %v", eventoFechaID, err)
		return false, err
	}
	return estado == 1, nil
}

// VerificarEventoFechaPerteneceAEvento: true si el evento_fecha pertenece al evento.
func (r *EventoFecha) VerificarEventoFechaPerteneceAEvento(eventoFechaID, eventoID int64) (bool, error) {
	var count int64
	res := r.PostgresqlDB.
		Table("evento_fecha").
		Where("evento_fecha_id = ? AND evento_id = ?", eventoFechaID, eventoID).
		Count(&count)
	if res.Error != nil {
		r.logger.Errorf("VerificarEventoFechaPerteneceAEvento(ef=%d,e=%d): %v", eventoFechaID, eventoID, res.Error)
		return false, res.Error
	}
	return count == 1, nil
}

// ContarEventoFechasActivas: total de fechas activas de un evento (útil para UI/validaciones).
func (r *EventoFecha) ContarEventoFechasActivas(eventoID int64) (int64, error) {
	var count int64
	res := r.PostgresqlDB.
		Table("evento_fecha").
		Where("evento_id = ? AND estado = 1", eventoID).
		Count(&count)
	if res.Error != nil {
		r.logger.Errorf("ContarEventoFechasActivas(%d): %v", eventoID, res.Error)
		return 0, res.Error
	}
	return count, nil
}

// ActivarEventoFecha cambia el estado de un evento_fecha a activo (1)
func (r *EventoFecha) ActivarEventoFecha(eventoFechaID int64) error {
	res := r.PostgresqlDB.
		Table("evento_fecha").
		Where("evento_fecha_id = ?", eventoFechaID).
		Update("estado", int16(1))
	if res.Error != nil {
		r.logger.Errorf("ActivarEventoFecha(%d): %v", eventoFechaID, res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DesactivarEventoFecha cambia el estado de un evento_fecha a inactivo (0)
func (r *EventoFecha) DesactivarEventoFecha(eventoFechaID int64) error {
	res := r.PostgresqlDB.
		Table("evento_fecha").
		Where("evento_fecha_id = ?", eventoFechaID).
		Update("estado", int16(0))
	if res.Error != nil {
		r.logger.Errorf("DesactivarEventoFecha(%d): %v", eventoFechaID, res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
