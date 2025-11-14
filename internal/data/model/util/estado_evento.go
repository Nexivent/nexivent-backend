package util

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type EstadoEvento int16

const (
	EventoBorrador EstadoEvento = iota
	EventoPublicado
	EventoCancelado
)

var ErrInvalidEstadoEvento = errors.New("estado de evento inválido")

// == equivalente a getCodigo() ==
func (e EstadoEvento) Codigo() int16 { return int16(e) }

// == equivalente a valueOfCodigo(Integer codigo) ==
func ValueOfEstadoEventoCodigo(c int16) (EstadoEvento, error) {
	switch c {
	case 0:
		return EventoBorrador, nil
	case 1:
		return EventoPublicado, nil
	case 2:
		return EventoCancelado, nil
	default:
		return 0, fmt.Errorf("código de estado inválido: %d", c)
	}
}

func (e EstadoEvento) String() string {
	switch e {
	case EventoBorrador:
		return "BORRADOR"
	case EventoPublicado:
		return "PUBLICADO"
	case EventoCancelado:
		return "CANCELADO"
	default:
		return "DESCONOCIDO"
	}
}

func (e EstadoEvento) IsValid() bool { return e >= EventoBorrador && e <= EventoCancelado }

/* ---- Integración con database/sql (columna SMALLINT) ---- */

func (e EstadoEvento) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, fmt.Errorf("estado inválido: %d", e)
	}
	// driver.Value acepta int64; Postgres SMALLINT castea sin problema
	return int64(e), nil
}

func (e *EstadoEvento) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		*e = EstadoEvento(v)
	case int32:
		*e = EstadoEvento(v)
	case int16:
		*e = EstadoEvento(v)
	case []byte:
		var n int16
		if _, err := fmt.Sscanf(string(v), "%d", &n); err != nil {
			return fmt.Errorf("scan EstadoEvento: %w", err)
		}
		*e = EstadoEvento(n)
	case string:
		var n int16
		if _, err := fmt.Sscanf(v, "%d", &n); err != nil {
			return fmt.Errorf("scan EstadoEvento: %w", err)
		}
		*e = EstadoEvento(n)
	default:
		return fmt.Errorf("tipo no soportado para EstadoEvento: %T", src)
	}
	if !e.IsValid() {
		return fmt.Errorf("estado inválido: %d", *e)
	}
	return nil
}

func (e *EstadoEvento) UnmarshalJSON(data []byte) error {
	unquotedData, err := strconv.Unquote(string(data))
	if err != nil {
		return ErrInvalidEstadoEvento
	}

	// Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedData, " ")

	// Sanity check the parts of the string to make sure it was in the expected format.
	// If it isn't, we return the ErrInvalidEstadoEvento error again.
	if len(parts) > 1 {
		return ErrInvalidEstadoEvento
	}

	switch strings.ToUpper(parts[0]) {
	case "BORRADOR":
		*e = EventoBorrador
	case "PUBLICADO":
		*e = EventoPublicado
	case "CANCELADO":
		*e = EventoCancelado
	default:
		return ErrInvalidEstadoEvento
	}

	return nil
}

func (e EstadoEvento) MarshalJSON() ([]byte, error) {
	return fmt.Appendf([]byte{}, `"%s"`, e.String()), nil
}
