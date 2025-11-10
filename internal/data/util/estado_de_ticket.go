package util

import (
	"database/sql/driver"
	"fmt"
)

type EstadoDeTicket int16

const (
	TicketDisponible EstadoDeTicket = iota // 0
	TicketVendido                          // 1
	TicketUsado                            // 2
	TicketCancelado                        // 3
)

func (e EstadoDeTicket) Codigo() int16 { return int16(e) }

func ValueOfEstadoDeTicketCodigo(c int16) (EstadoDeTicket, error) {
	switch c {
	case 0:
		return TicketDisponible, nil
	case 1:
		return TicketVendido, nil
	case 2:
		return TicketUsado, nil
	case 3:
		return TicketCancelado, nil
	default:
		return 0, fmt.Errorf("código de estado de ticket inválido: %d", c)
	}
}

func (e EstadoDeTicket) String() string {
	switch e {
	case TicketDisponible:
		return "DISPONIBLE"
	case TicketVendido:
		return "VENDIDO"
	case TicketUsado:
		return "USADO"
	case TicketCancelado:
		return "CANCELADO"
	default:
		return "DESCONOCIDO"
	}
}

func (e EstadoDeTicket) IsValid() bool {
	return e >= TicketDisponible && e <= TicketCancelado
}

/* ---- Integración con database/sql (columna SMALLINT) ---- */

func (e EstadoDeTicket) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, fmt.Errorf("estado de ticket inválido: %d", e)
	}
	return int64(e), nil // el driver lo castea a SMALLINT en Postgres
}

func (e *EstadoDeTicket) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		*e = EstadoDeTicket(v)
	case int32:
		*e = EstadoDeTicket(v)
	case int16:
		*e = EstadoDeTicket(v)
	case []byte:
		var n int16
		if _, err := fmt.Sscanf(string(v), "%d", &n); err != nil {
			return fmt.Errorf("scan EstadoDeTicket: %w", err)
		}
		*e = EstadoDeTicket(n)
	case string:
		var n int16
		if _, err := fmt.Sscanf(v, "%d", &n); err != nil {
			return fmt.Errorf("scan EstadoDeTicket: %w", err)
		}
		*e = EstadoDeTicket(n)
	default:
		return fmt.Errorf("tipo no soportado para EstadoDeTicket: %T", src)
	}
	if !e.IsValid() {
		return fmt.Errorf("estado de ticket inválido: %d", *e)
	}
	return nil
}
