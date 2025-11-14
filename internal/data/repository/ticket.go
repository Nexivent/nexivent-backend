package repository

import (
	"errors"

	"github.com/Nexivent/nexivent-backend/internal/data/model"
	util "github.com/Nexivent/nexivent-backend/internal/data/model/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Ticket struct {
	DB *gorm.DB
}

func NewTicketsController(
	postgresqlDB *gorm.DB,
) *Ticket {
	return &Ticket{
		DB: postgresqlDB,
	}
}

// VerificarOrdenConfirmada: true si la orden existe y está CONFIRMADA (enum).
func (c *Ticket) VerificarOrdenConfirmada(orderID int64) (bool, error) {
	var estInt int16
	err := c.DB.
		Table("orden_de_compra").
		Select("estado_de_orden").
		Where("orden_de_compra_id = ?", orderID).
		Row().
		Scan(&estInt)

	if err != nil {
		return false, err
	}
	return util.EstadoOrden(estInt) == util.OrdenConfirmada, nil
}

// VerificarTicketsExistentes: true si la orden ya tiene tickets creados.
func (c *Ticket) VerificarTicketsExistentes(orderID int64) (bool, error) {
	var count int64
	res := c.DB.
		Table("ticket").
		Where("orden_de_compra_id = ?", orderID).
		Count(&count)
	if res.Error != nil {
		return false, res.Error
	}
	return count > 0, nil
}

// VerificarStockDisponible: revisa capacidad del sector asociado a la tarifa.
func (c *Ticket) VerificarStockDisponible(tarifaID int64, cantidad int64) (bool, error) {
	if tarifaID <= 0 || cantidad <= 0 {
		return false, gorm.ErrInvalidData
	}

	var total, vendidas int64
	row := c.DB.
		Table("sector s").
		Select("s.total_entradas, s.cant_vendidas").
		Joins("JOIN tarifa t ON t.sector_id = s.sector_id").
		Where("t.tarifa_id = ?", tarifaID).
		Row()

	if err := row.Scan(&total, &vendidas); err != nil {
		return false, err
	}

	return vendidas+cantidad <= total, nil
}

// VerificarTicketPerteneceAOrden: true si ticket_id pertenece a la orden dada.
func (c *Ticket) VerificarTicketPerteneceAOrden(ticketID, orderID int64) (bool, error) {
	var count int64
	res := c.DB.
		Table("ticket").
		Where("ticket_id = ? AND orden_de_compra_id = ?", ticketID, orderID).
		Count(&count)
	if res.Error != nil {
		return false, res.Error
	}
	return count == 1, nil
}

// CrearTicket: inserta un solo ticket (el BO define QR y estado).
func (c *Ticket) CrearTicket(ticket *model.Ticket) error {
	if ticket == nil {
		return gorm.ErrInvalidData
	}
	if err := c.DB.Create(ticket).Error; err != nil {
		return err
	}
	return nil
}

// CrearTicketsBatch: inserta varios tickets (BO validó todo antes).
func (c *Ticket) CrearTicketsBatch(tickets []model.Ticket) error {
	if len(tickets) == 0 {
		return nil
	}
	if err := c.DB.Create(&tickets).Error; err != nil {
		return err
	}
	return nil
}

// IncrementarVendidasPorSector: suma cantidad a cant_vendidas (sector).
func (c *Ticket) IncrementarVendidasPorSector(sectorID int64, cantidad int64) error {
	if sectorID <= 0 || cantidad <= 0 {
		return gorm.ErrInvalidData
	}
	res := c.DB.
		Table("sector").
		Where("sector_id = ?", sectorID).
		UpdateColumn("cant_vendidas", gorm.Expr("cant_vendidas + ?", cantidad))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ObtenerSectorPorTarifa: retorna sector_id para una tarifa.
func (c *Ticket) ObtenerSectorPorTarifa(tarifaID int64) (int64, error) {
	var sectorID int64
	row := c.DB.
		Table("tarifa").
		Select("sector_id").
		Where("tarifa_id = ?", tarifaID).
		Row()
	if err := row.Scan(&sectorID); err != nil {
		return 0, err
	}
	return sectorID, nil
}

// CambiarEstadoTicket: actualiza el estado de un ticket usando tu enum model.EstadoDeTicket.
func (c *Ticket) CambiarEstadoTicket(ticketID int64, nuevo util.EstadoDeTicket) error {
	if !nuevo.IsValid() {
		return errors.New("estado de ticket inválido")
	}
	res := c.DB.
		Table("ticket").
		Where("ticket_id = ?", ticketID).
		Update("estado_de_ticket", nuevo.Codigo())
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// CambiarEstadoTicketsDeOrden: actualiza el estado de todos los tickets de una orden.
func (c *Ticket) CambiarEstadoTicketsDeOrden(orderID int64, nuevo util.EstadoDeTicket) (int64, error) {
	if !nuevo.IsValid() {
		return 0, errors.New("estado de ticket inválido")
	}
	res := c.DB.
		Table("ticket").
		Where("orden_de_compra_id = ?", orderID).
		Update("estado_de_ticket", nuevo.Codigo())
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

// ObtenerTicketsPorOrden: devuelve tickets (modelo puro) de una orden (para el BO).
func (c *Ticket) ObtenerTicketsPorOrden(orderID int64) ([]model.Ticket, error) {
	var ts []model.Ticket
	if err := c.DB.
		Where("orden_de_compra_id = ?", orderID).
		Find(&ts).Error; err != nil {
		return nil, err
	}
	return ts, nil
}

// LockTicketsDeOrden: bloquea filas de tickets por orden (si el BO necesita transacción estricta).
func (c *Ticket) LockTicketsDeOrden(orderID int64) ([]model.Ticket, error) {
	var ts []model.Ticket
	if err := c.DB.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("orden_de_compra_id = ?", orderID).
		Find(&ts).Error; err != nil {
		return nil, err
	}
	return ts, nil
}
