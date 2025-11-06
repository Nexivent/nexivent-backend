package model

import (
	"database/sql/driver"
	"fmt"
)

type EstadoNotificacion int16

const (
	NotificacionEnviada   EstadoNotificacion = iota // 0
	NotificacionNoEnviada                           // 1
)

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
