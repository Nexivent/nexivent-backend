package util

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type EstadoOrden int16

const (
	OrdenTemporal EstadoOrden = iota
	OrdenConfirmada
	OrdenCancelada
)

var ErrInvalidEstadoOrden = errors.New("estado de orden de compra inválido")

func (e EstadoOrden) Codigo() int16 { return int16(e) }

func (e EstadoOrden) String() string {
	switch e {
	case OrdenTemporal:
		return "TEMPORAL"
	case OrdenConfirmada:
		return "CONFIRMADA"
	case OrdenCancelada:
		return "CANCELADA"
	default:
		return "DESCONOCIDO"
	}
}

func (e EstadoOrden) IsValid() bool {
	return e == OrdenTemporal || e == OrdenConfirmada || e == OrdenCancelada
}

/* ---- Integración con database/sql (columna SMALLINT) ---- */

func (e EstadoOrden) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, fmt.Errorf("estado_de_orden inválido: %d", e)
	}
	return int64(e), nil
}

func (e *EstadoOrden) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		*e = EstadoOrden(v)
	case int32:
		*e = EstadoOrden(v)
	case int16:
		*e = EstadoOrden(v)
	case []byte:
		var n int16
		if _, err := fmt.Sscanf(string(v), "%d", &n); err != nil {
			return fmt.Errorf("scan EstadoOrden: %w", err)
		}
		*e = EstadoOrden(n)
	case string:
		var n int16
		if _, err := fmt.Sscanf(v, "%d", &n); err != nil {
			return fmt.Errorf("scan EstadoOrden: %w", err)
		}
		*e = EstadoOrden(n)
	default:
		return fmt.Errorf("tipo no soportado para EstadoOrden: %T", src)
	}
	if !e.IsValid() {
		return fmt.Errorf("estado_de_orden inválido: %d", *e)
	}
	return nil
}

func (e *EstadoOrden) UnmarshalJSON(data []byte) error {
	unquotedData, err := strconv.Unquote(string(data))
	if err != nil {
		return ErrInvalidEstadoOrden
	}

	// Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedData, " ")

	// Sanity check the parts of the string to make sure it was in the expected format.
	// If it isn't, we return the ErrInvalidEstadoOrden error again.
	if len(parts) > 1 {
		return ErrInvalidEstadoOrden
	}

	switch strings.ToUpper(parts[0]) {
	case "TEMPORAL":
		*e = OrdenTemporal
	case "CONFIRMADA":
		*e = OrdenConfirmada
	case "CANCELADA":
		*e = OrdenCancelada
	default:
		return ErrInvalidEstadoOrden
	}

	return nil
}

func (e EstadoOrden) MarshalJSON() ([]byte, error) {
	return fmt.Appendf([]byte{}, `"%s"`, e.String()), nil
}
