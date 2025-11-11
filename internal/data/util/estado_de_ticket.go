package util

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type EstadoDeTicket int16

const (
	TicketDisponible EstadoDeTicket = iota
	TicketVendido
	TicketUsado
	TicketCancelado
)

var ErrInvalidEstadoTicket = errors.New("estado de ticket inválido")

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

func (e *EstadoDeTicket) UnmarshalJSON(data []byte) error {
	unquotedData, err := strconv.Unquote(string(data))
	if err != nil {
		return ErrInvalidEstadoTicket
	}

	// Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedData, " ")

	// Sanity check the parts of the string to make sure it was in the expected format.
	// If it isn't, we return the ErrInvalidEstadoTicket error again.
	if len(parts) > 1 {
		return ErrInvalidEstadoTicket
	}

	switch strings.ToUpper(parts[0]) {
	case "DISPONIBLE":
		*e = TicketDisponible
	case "VENDIDO":
		*e = TicketVendido
	case "USADO":
		*e = TicketUsado
	case "CANCELADO":
		*e = TicketCancelado
	default:
		return ErrInvalidEstadoTicket
	}

	return nil
}

func (e EstadoDeTicket) MarshalJSON() ([]byte, error) {
	return fmt.Appendf([]byte{}, `"%s"`, e.String()), nil
}
