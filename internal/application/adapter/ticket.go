package adapter

import (
	"fmt"
	"time"  // ← NECESARIO para time.Now()
	
	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/dao/model"
	util "github.com/Nexivent/nexivent-backend/internal/dao/model/util"  
	daoPostgresql "github.com/Nexivent/nexivent-backend/internal/dao/repository"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
)

type Ticket struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.NexiventPsqlEntidades
}

func NewTicketAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.NexiventPsqlEntidades,
) *Ticket {
	return &Ticket{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Genera un código QR pseudo-único (puedes ajustar el formato)
func generarCodigoQR(orderID int64, tarifaID int64, corr int64) string {
	return fmt.Sprintf("TCK-%d-%d-%d", orderID, tarifaID, corr)
}

func (t *Ticket) EmitirTickets(orderID int64) (*schemas.TicketIssueResponse, *errors.Error) {
	if orderID <= 0 {
		return nil, &errors.UnprocessableEntityError.InvalidReservationId
	}

	daoTicket := t.DaoPostgresql.Ticket

	// 1) Orden debe existir y estar CONFIRMADA
	ok, err := daoTicket.VerificarOrdenConfirmada(orderID)
	if err != nil {
		t.logger.Errorf("EmitirTickets.VerificarOrdenConfirmada(%d): %v", orderID, err)
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}
		return nil, &errors.InternalServerError.Default
	}
	if !ok {
		// Orden existe pero no está confirmada / no válida
		return nil, &errors.ObjectNotFoundError.EventoNotFound
	}

	// 2) Verificar que aún no existan tickets
	yaTiene, err := daoTicket.VerificarTicketsExistentes(orderID)
	if err != nil {
		t.logger.Errorf("EmitirTickets.VerificarTicketsExistentes(%d): %v", orderID, err)
		return nil, &errors.InternalServerError.Default
	}
	if yaTiene {
		// Orden ya procesada
		return nil, &errors.ObjectNotFoundError.EventoNotFound
	}

	// 3) Obtener detalles de la orden
	detalles, err := daoTicket.ObtenerDetallesOrden(orderID)
	if err != nil {
		t.logger.Errorf("EmitirTickets.ObtenerDetallesOrden(%d): %v", orderID, err)
		return nil, &errors.InternalServerError.Default
	}
	if len(detalles) == 0 {
		return nil, &errors.ObjectNotFoundError.EventoNotFound
	}

	// 4) Verificar stock por cada detalle
	for _, d := range detalles {
		ok, err := daoTicket.VerificarStockDisponible(d.TarifaID, d.Cantidad)
		if err != nil {
			t.logger.Errorf("EmitirTickets.VerificarStockDisponible(order=%d, tarifa=%d): %v", orderID, d.TarifaID, err)
			return nil, &errors.InternalServerError.Default
		}
		if !ok {
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}
	}

	// 5) Crear modelos de tickets (solo lógica de negocio, sin BD)
	var ticketsAInsertar []model.Ticket
	for _, d := range detalles {
		for i := int64(0); i < d.Cantidad; i++ {
			qr := generarCodigoQR(orderID, d.TarifaID, i)
			ticketModel := model.Ticket{
				OrdenDeCompraID: &orderID,
				EventoFechaID:   d.EventoFechaID,
				TarifaID:        d.TarifaID,
				CodigoQR:        qr,
				EstadoDeTicket:  util.EstadoDeTicket(0).Codigo(), // DISPONIBLE
			}
			ticketsAInsertar = append(ticketsAInsertar, ticketModel)
		}
	}

	// 6) Insertar tickets en BD (DAO)
	if err := daoTicket.CrearTicketsBatch(ticketsAInsertar); err != nil {
		t.logger.Errorf("EmitirTickets.CrearTicketsBatch(order=%d): %v", orderID, err)
		return nil, &errors.InternalServerError.Default
	}

	// 7) Actualizar cant_vendidas por sector
	for _, d := range detalles {
		sectorID, err := daoTicket.ObtenerSectorPorTarifa(d.TarifaID)
		if err != nil {
			t.logger.Errorf("EmitirTickets.ObtenerSectorPorTarifa(tarifa=%d): %v", d.TarifaID, err)
			return nil, &errors.InternalServerError.Default
		}
		if err := daoTicket.IncrementarVendidasPorSector(sectorID, d.Cantidad); err != nil {
			t.logger.Errorf("EmitirTickets.IncrementarVendidasPorSector(sector=%d): %v", sectorID, err)
			return nil, &errors.InternalServerError.Default
		}
	}

	// 8) Traer info para respuesta (JOINs complejos ya encapsulados en DAO)
	infoRows, err := daoTicket.ObtenerTicketsInfoPorOrden(orderID)
	if err != nil {
		t.logger.Errorf("EmitirTickets.ObtenerTicketsInfoPorOrden(order=%d): %v", orderID, err)
		return nil, &errors.InternalServerError.Default
	}

	// 9) Mapear a schemas
	ticketsResp := make([]schemas.TicketEmitido, 0, len(infoRows))
	for _, row := range infoRows {
		ticketsResp = append(ticketsResp, schemas.TicketEmitido{
			IdTicket:     row.ID,
			CodigoQR:     row.CodigoQR,
			Estado:       util.EstadoDeTicket(row.Estado).String(),
			TituloEvento: row.Titulo,
			FechaEvento:  row.FechaEvento.Format("2006-01-02"),
			HoraInicio:   row.HoraInicio.Format("15:04"),
			Sector:       row.SectorTipo,
		})
	}

	return &schemas.TicketIssueResponse{
		Tickets: ticketsResp,
	}, nil
}

func (t *Ticket) CancelarTickets(req *schemas.TicketCancelRequest) (*schemas.TicketCancelResponse, *errors.Error) {
	if len(req.IdTickets) == 0 {
		return nil, &errors.UnprocessableEntityError.InvalidReservationId
	}

	daoTicket := t.DaoPostgresql.Ticket

	// 1) Obtener estado y tarifa de los tickets
	rows, err := daoTicket.ObtenerTicketsEstadoTarifaPorIDs(req.IdTickets)
	if err != nil {
		t.logger.Errorf("CancelarTickets.ObtenerTicketsEstadoTarifaPorIDs: %v", err)
		return nil, &errors.InternalServerError.Default
	}

	found := make(map[int64]daoPostgresql.TicketEstadoTarifa, len(rows))
	for _, r := range rows {
		found[r.ID] = r
	}

	cancelados := []schemas.TicketCancelado{}
	noEncontrados := []int64{}
	noCancelables := []int64{}

	for _, id := range req.IdTickets {
		row, ok := found[id]
		if !ok {
			noEncontrados = append(noEncontrados, id)
			continue
		}

		// No se pueden cancelar USADO(1) ni CANCELADO(2)
		if row.EstadoDeTicket == 1 || row.EstadoDeTicket == 2 {
			noCancelables = append(noCancelables, id)
			continue
		}

		// Cambiar a CANCELADO
		if err := daoTicket.CambiarEstadoTicket(id, util.EstadoDeTicket(2)); err != nil {
			t.logger.Errorf("CancelarTickets.CambiarEstadoTicket(%d): %v", id, err)
			noCancelables = append(noCancelables, id)
			continue
		}

		// Sector del ticket (vía tarifa)
		sectorID, err := daoTicket.ObtenerSectorPorTarifa(row.TarifaID)
		if err != nil {
			t.logger.Errorf("CancelarTickets.ObtenerSectorPorTarifa(tarifa=%d): %v", row.TarifaID, err)
			noCancelables = append(noCancelables, id)
			continue
		}

		// Decrementar vendidas (sube disponibilidad)
		if err := daoTicket.DecrementarVendidasPorSector(sectorID, 1); err != nil {
			t.logger.Errorf("CancelarTickets.DecrementarVendidasPorSector(sector=%d): %v", sectorID, err)
			noCancelables = append(noCancelables, id)
			continue
		}

		cancelados = append(cancelados, schemas.TicketCancelado{
			IdTicket: id,
			Estado:   util.EstadoDeTicket(3).String(),
		})
	}

	if len(cancelados) == 0 {
		// “Error al cancelar” según contrato
		return nil, &errors.ObjectNotFoundError.EventoNotFound
	}

	resp := &schemas.TicketCancelResponse{
		Cancelados:    cancelados,
		NoEncontrados: noEncontrados,
		NoCancelables: noCancelables,
		Mensaje:       "Tickets cancelados correctamente.",
	}
	return resp, nil
}

func (t *Ticket) EmitirTicketsConInfo(
	req *schemas.EmitirTicketsRequest,
) (*schemas.EmitirTicketsResponse, *errors.Error) {

	if req.OrderID == 0 || req.UserID == 0 || len(req.Tickets) == 0 {
		return nil, &errors.UnprocessableEntityError.InvalidRequestBody
	}

	orden, err := t.DaoPostgresql.OrdenDeCompra.ObtenerOrdenBasica(req.OrderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			t.logger.Warnf("Orden %d no encontrada", req.OrderID)
			return nil, &errors.ObjectNotFoundError.EventoNotFound
		}
		t.logger.Errorf("EmitirTicketsConInfo.ObtenerOrden(%d): %v", req.OrderID, err)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	if orden.EstadoDeOrden != util.OrdenConfirmada.Codigo() {
		t.logger.Warnf("Orden %d no está confirmada (estado: %d)", req.OrderID, orden.EstadoDeOrden)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	if orden.UsuarioID != req.UserID {
		t.logger.Warnf("Usuario %d no es dueño de la orden %d", req.UserID, req.OrderID)
		return nil, &errors.BadRequestError.EventoNotFound
	}

	var ticketsGenerados []schemas.TicketGenerado

	for _, ticketInfo := range req.Tickets {
		for i := 0; i < ticketInfo.Cantidad; i++ {
			timestamp := time.Now().UnixNano()
			codigoQR := fmt.Sprintf("QR-%d-%d-%d-%d", timestamp, req.OrderID, ticketInfo.IdTarifa, i)

			ordenID := req.OrderID
			ticket := &model.Ticket{
				OrdenDeCompraID: &ordenID,
				EventoFechaID:   req.IdFechaEvento,
				TarifaID:        ticketInfo.IdTarifa,
				CodigoQR:        codigoQR,
				EstadoDeTicket:  int16(1), // 1 = VENDIDO
			}

			if err := t.DaoPostgresql.Ticket.Crear(ticket); err != nil {
				t.logger.Errorf("EmitirTicketsConInfo.CrearTicket: %v", err)
				return nil, &errors.BadRequestError.EventoNotCreated
			}

			ticketsGenerados = append(ticketsGenerados, schemas.TicketGenerado{
				IdTicket: fmt.Sprintf("%d", ticket.ID),
				CodigoQR: ticket.CodigoQR,
				Estado:   "VENDIDO",
				Zona:     ticketInfo.NombreZona,
			})
		}
	}

	t.logger.Infof("✅ Tickets generados para orden %d: %d tickets", req.OrderID, len(ticketsGenerados))

	resp := &schemas.EmitirTicketsResponse{
		Tickets: ticketsGenerados,
		OrderID: req.OrderID,
	}
	return resp, nil
}
