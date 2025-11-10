package data

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/validator"
	"github.com/jmoiron/sqlx"
)

type Evento struct {
	ID                  int64      `json:"id"`
	Nombre              string      `json:"nombre"`
	Organizador         Organizador `json:"organizador"`
	FechaCreacion       time.Time   `json:"-"`
	FechaRealizacion    time.Time   `json:"fecha_realizacion"`
	CantidadAsientos    uint32      `json:"cantidad_asientos"`
	AsientosDisponibles uint32      `json:"asientos_disponibles"`
	AsientosOcupados    uint32      `json:"asientos_ocupados"`
}

type EventoModel struct {
	DB *sqlx.DB
}

func ValidateEvento(v *validator.Validator, evento *Evento) {
	// Validar nombre del evento
	v.Check(evento.Nombre != "", "nombre", "debe ser proporcionado")
	v.Check(len(evento.Nombre) <= 500, "nombre", "no debe tener más de 500 caracteres")

	// Validar fecha de realización
	v.Check(!evento.FechaRealizacion.IsZero(), "fecha_realizacion", "debe ser proporcionada")
	v.Check(evento.FechaRealizacion.After(time.Now()), "fecha_realizacion", "debe ser una fecha futura")

	// Validar cantidad de asientos
	v.Check(evento.CantidadAsientos > 0, "cantidad_asientos", "debe ser mayor a 0")
	v.Check(evento.CantidadAsientos <= 1000000, "cantidad_asientos", "no debe exceder 1,000,000")

	// Validar coherencia entre asientos
	v.Check(evento.AsientosDisponibles <= evento.CantidadAsientos, "asientos_disponibles", "no puede exceder la cantidad total de asientos")
	v.Check(evento.AsientosOcupados <= evento.CantidadAsientos, "asientos_ocupados", "no puede exceder la cantidad total de asientos")
	v.Check(evento.AsientosDisponibles+evento.AsientosOcupados == evento.CantidadAsientos, "asientos", "la suma de asientos disponibles y ocupados debe igual a la cantidad total")
}

func (m EventoModel) Insert(evento *Evento) error {
	query := `
		INSERT INTO eventos (nombre, organizador_ruc, fecha_realizacion, cantidad_asientos, asientos_disponibles, asientos_ocupados)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, fecha_creacion`

	args := []any{
		evento.Nombre,
		evento.Organizador.RUC,
		evento.FechaRealizacion,
		evento.CantidadAsientos,
		evento.AsientosDisponibles,
		evento.AsientosOcupados,
	}

	err := m.DB.QueryRow(query, args...).Scan(&evento.ID, &evento.FechaCreacion)
	if err != nil {
		return err
	}

	return nil
}

func (m EventoModel) Get(id int64) (*Evento, error) {
	return nil, nil
}

func (m EventoModel) Update(evento *Evento) error {
	return nil
}

func (m EventoModel) Delete(id int64) error {
	return nil
}
