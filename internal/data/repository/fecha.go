package repository

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/model"
	"gorm.io/gorm"
)

type Fecha struct {
	DB *gorm.DB
}

func (f *Fecha) CrearFecha(Fecha *model.Fecha) error {
	respuesta := f.DB.Create(Fecha)

	if respuesta.Error != nil {
		return respuesta.Error
	}

	return nil
}

func (f *Fecha) ObtenerPorFecha(Fecha time.Time) (*model.Fecha, error) {
	var fechaObtenida *model.Fecha
	respuesta := f.DB.Where("fecha_evento = ?", Fecha).Find(&fechaObtenida)

	if respuesta.Error != nil {
		return nil, respuesta.Error
	}

	return fechaObtenida, nil
}
