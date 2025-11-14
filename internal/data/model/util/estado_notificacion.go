package util

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type EstadoNotificacion int16

const (
	NotificacionEnviada EstadoNotificacion = iota
	NotificacionNoEnviada
)

var ErrInvalidEstadoNotificacion = errors.New("estado de notificación inválido")

func (e EstadoNotificacion) Codigo() int16 { return int16(e) }

func ValueOfEstadoNotificacionCodigo(c int16) (EstadoNotificacion, error) {
	switch c {
	case 0:
		return NotificacionEnviada, nil
	case 1:
		return NotificacionNoEnviada, nil
	default:
		return 0, fmt.Errorf("código de estado de notificación inválido: %d", c)
	}
}

func (e EstadoNotificacion) String() string {
	switch e {
	case NotificacionEnviada:
		return "ENVIADO"
	case NotificacionNoEnviada:
		return "NO_ENVIADO"
	default:
		return "DESCONOCIDO"
	}
}

func (e EstadoNotificacion) IsValid() bool {
	return e == NotificacionEnviada || e == NotificacionNoEnviada
}

/* ---- Integración con database/sql (columna SMALLINT) ---- */

func (e EstadoNotificacion) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, fmt.Errorf("estado de notificación inválido: %d", e)
	}
	return int64(e), nil
}

func (e *EstadoNotificacion) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		*e = EstadoNotificacion(v)
	case int32:
		*e = EstadoNotificacion(v)
	case int16:
		*e = EstadoNotificacion(v)
	case []byte:
		var n int16
		if _, err := fmt.Sscanf(string(v), "%d", &n); err != nil {
			return fmt.Errorf("scan EstadoNotificacion: %w", err)
		}
		*e = EstadoNotificacion(n)
	case string:
		var n int16
		if _, err := fmt.Sscanf(v, "%d", &n); err != nil {
			return fmt.Errorf("scan EstadoNotificacion: %w", err)
		}
		*e = EstadoNotificacion(n)
	default:
		return fmt.Errorf("tipo no soportado para EstadoNotificacion: %T", src)
	}
	if !e.IsValid() {
		return fmt.Errorf("estado de notificación inválido: %d", *e)
	}
	return nil
}

func (e *EstadoNotificacion) UnmarshalJSON(data []byte) error {
	unquotedData, err := strconv.Unquote(string(data))
	if err != nil {
		return ErrInvalidEstadoNotificacion
	}

	// Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedData, " ")

	// Sanity check the parts of the string to make sure it was in the expected format.
	// If it isn't, we return the ErrInvalidEstadoNotificacion error again.
	if len(parts) > 1 {
		return ErrInvalidEstadoNotificacion
	}

	switch strings.ToUpper(parts[0]) {
	case "ENVIADO":
		*e = NotificacionEnviada
	case "NO_ENVIADO":
		*e = NotificacionNoEnviada
	default:
		return ErrInvalidEstadoNotificacion
	}

	return nil
}

func (e EstadoNotificacion) MarshalJSON() ([]byte, error) {
	return fmt.Appendf([]byte{}, `"%s"`, e.String()), nil
}
