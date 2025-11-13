package repository

import (
	"github.com/Nexivent/nexivent-backend/internal/data/model"
	"gorm.io/gorm"
)

type TipoDeTicket struct {
	DB *gorm.DB
}

func (t *TipoDeTicket) CrearTipoDeTicket(TipoDeTicket *model.TipoDeTicket) error {
	resultado := t.DB.Create(TipoDeTicket)

	if resultado.Error != nil {
		return resultado.Error
	}

	return nil
}

func (t *TipoDeTicket) ActualizarTipoDeTicketr(TipoDeTicket *model.TipoDeTicket) error {
	respuesta := t.DB.Save(TipoDeTicket)

	if respuesta.Error != nil {
		return respuesta.Error
	}

	return nil
}

func (t *TipoDeTicket) ListarTipoTicketPorEventoID(eventoID uint64) ([]model.TipoDeTicket, error) {
	var list []model.TipoDeTicket
	if err := t.DB.
		Where("evento_id = ?", eventoID).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
