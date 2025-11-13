package repository

import (
	"github.com/Loui27/nexivent-backend/internal/dao/model"
	"github.com/Loui27/nexivent-backend/logging"
	"gorm.io/gorm"
)

type TipoDeTicket struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewTipoDeTicketController(
	logger logging.Logger,
	postgreesqlDB *gorm.DB,
) *TipoDeTicket {
	return &TipoDeTicket{
		logger:       logger,
		PostgresqlDB: postgreesqlDB,
	}
}

func (t *TipoDeTicket) CrearTipoDeTicket(TipoDeTicket *model.TipoDeTicket) error {
	resultado := t.PostgresqlDB.Create(TipoDeTicket)

	if resultado.Error != nil {
		return resultado.Error
	}

	return nil
}

func (t *TipoDeTicket) ActualizarTipoDeTicketr(TipoDeTicket *model.TipoDeTicket) error {
	respuesta := t.PostgresqlDB.Save(TipoDeTicket)

	if respuesta.Error != nil {
		return respuesta.Error
	}

	return nil
}

// ListarTipoTicketPorEventoID: devuelve TODAS las filas por evento_id (sin filtrar estado)
func (t *TipoDeTicket) ListarTipoTicketPorEventoID(eventoID int64) ([]model.TipoDeTicket, error) {
	var list []model.TipoDeTicket
	if err := t.PostgresqlDB.
		Where("evento_id = ?", eventoID).
		Find(&list).Error; err != nil {
		t.logger.Errorf("ListarTipoTicketPorEventoID(%d): %v", eventoID, err)
		return nil, err
	}
	return list, nil
}
