package repository

import (
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
