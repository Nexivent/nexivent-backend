package controller

import (
	"github.com/Nexivent/nexivent-backend/errors"
	"github.com/Nexivent/nexivent-backend/internal/application/adapter"
	"github.com/Nexivent/nexivent-backend/internal/schemas"
	"github.com/Nexivent/nexivent-backend/logging"
)

type TicketController struct {
	Logger        logging.Logger
	TicketAdapter *adapter.Ticket
}

func NewTicketController(
	logger logging.Logger,
	ticketAdapter *adapter.Ticket,
) *TicketController {
	return &TicketController{
		Logger:        logger,
		TicketAdapter: ticketAdapter,
	}
}

func (tc *TicketController) EmitirTickets(orderID int64) (*schemas.TicketIssueResponse, *errors.Error) {
	return tc.TicketAdapter.EmitirTickets(orderID)
}

func (tc *TicketController) CancelarTickets(req schemas.TicketCancelRequest) (*schemas.TicketCancelResponse, *errors.Error) {
	return tc.TicketAdapter.CancelarTickets(&req)
}
