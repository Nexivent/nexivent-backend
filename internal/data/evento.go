// evento.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Evento struct {
	ID                  uuid.UUID         `db:"evento_id" json:"id"`
	OrganizadorID       uuid.UUID         `db:"organizador_id" json:"organizadorId"`
	CategoriaID         uuid.UUID         `db:"categoria_id" json:"categoriaId"`
	Titulo              string            `db:"titulo" json:"titulo"`
	Descripcion         string            `db:"descripcion" json:"descripcion"`
	Lugar               string            `db:"lugar" json:"lugar"`
	EventoEstado        util.EstadoEvento `db:"evento_estado" json:"eventoEstado"`
	CantMeGusta         int               `db:"cant_me_gusta" json:"cantMeGusta"`
	CantNoInteresa      int               `db:"cant_no_interesa" json:"cantNoInteresa"`
	CantVendidoTotal    int               `db:"cant_vendido_total" json:"cantVendidoTotal"`
	TotalRecaudado      float64           `db:"total_recaudado" json:"totalRecaudado"`
	Estado              util.Estado       `db:"estado" json:"estado"`
	UsuarioCreacion     *uuid.UUID        `db:"usuario_creacion" json:"usuarioCreacion"`
	FechaCreacion       time.Time         `db:"fecha_creacion" json:"fechaCreacion"`
	UsuarioModificacion *uuid.UUID        `db:"usuario_modificacion" json:"usuarioModificacion"`
	FechaModificacion   *time.Time        `db:"fecha_modificacion" json:"fechaModificacion"`
	ImagenDescripcion   string            `db:"imagen_descripcion" json:"imagenDescripcion"`
	ImagenPortada       string            `db:"imagen_portada" json:"imagenPortada"`
	VideoPresentacion   string            `db:"video_presentacion" json:"videoPresentacion"`
	ImagenEscenario     string            `db:"imagen_escenario" json:"imagenEscenario"`

	Organizador *Usuario
	Categoria   *Categoria

	Comentarios []Comentario
	Sectores    []Sector
	TiposTicket []TipoDeTicket
	Perfiles    []PerfilDePersona
	Fechas      []EventoFecha

	Cupones []Cupon
}

type EventoModel struct {
	DB *sqlx.DB
}

func ValidateEvento(v *validator.Validator, evento *Evento) {
	// Validar Titulo
	v.Check(evento.Titulo != "", "titulo", "el título es obligatorio")
	v.Check(len(evento.Titulo) <= 80, "titulo", "el título no debe exceder 80 caracteres")

	// Validar Descripcion
	v.Check(evento.Descripcion != "", "descripcion", "la descripción es obligatoria")

	// Validar Lugar
	v.Check(evento.Lugar != "", "lugar", "el lugar es obligatorio")
	v.Check(len(evento.Lugar) <= 80, "lugar", "el lugar no debe exceder 80 caracteres")

	// Validar IDs
	v.Check(evento.OrganizadorID != uuid.Nil, "organizador_id", "el ID del organizador es obligatorio")
	v.Check(evento.CategoriaID != uuid.Nil, "categoria_id", "el ID de categoría es obligatorio")

	// Validar contadores (no pueden ser negativos)
	v.Check(evento.CantMeGusta >= 0, "cant_me_gusta", "la cantidad de me gusta no puede ser negativa")
	v.Check(evento.CantNoInteresa >= 0, "cant_no_interesa", "la cantidad de no me interesa no puede ser negativa")
	v.Check(evento.CantVendidoTotal >= 0, "cant_vendido_total", "la cantidad vendida no puede ser negativa")
	v.Check(evento.TotalRecaudado >= 0, "total_recaudado", "el total recaudado no puede ser negativo")
}

func (e EventoModel) Insert(evento *Evento) error {
	query := `
		INSERT INTO evento (organizador_id, categoria_id, titulo, descripcion, lugar, 
			evento_estado, cant_me_gusta, cant_no_interesa, cant_vendido_total, 
			total_recaudado, estado, usuario_creacion)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING evento_id, fecha_creacion, fecha_modificacion`

	args := []any{
		evento.OrganizadorID,
		evento.CategoriaID,
		evento.Titulo,
		evento.Descripcion,
		evento.Lugar,
		evento.EventoEstado,
		evento.CantMeGusta,
		evento.CantNoInteresa,
		evento.CantVendidoTotal,
		evento.TotalRecaudado,
		evento.Estado,
		evento.UsuarioCreacion,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return e.DB.QueryRowContext(ctx, query, args...).Scan(
		&evento.ID,
		&evento.FechaCreacion,
		&evento.FechaModificacion,
	)
}

func (e EventoModel) Get(id uuid.UUID) (*Evento, error) {
	if id == uuid.Nil {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT evento_id, organizador_id, categoria_id, titulo, descripcion, lugar,
			evento_estado, cant_me_gusta, cant_no_interesa, cant_vendido_total,
			total_recaudado, estado, usuario_creacion, fecha_creacion,
			usuario_modificacion, fecha_modificacion
		FROM evento
		WHERE evento_id = $1 AND estado = 1`

	var evento Evento

	// Use the context.WithTimeout() function to create a context.Context which carries a
	// 3-second timeout deadline. Note that we're using the empty context.Background()
	// as the 'parent' context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Importantly, use defer to make sure that we cancel the context before the Get()
	// method returns.
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, id).Scan(
		&evento.ID,
		&evento.OrganizadorID,
		&evento.CategoriaID,
		&evento.Titulo,
		&evento.Descripcion,
		&evento.Lugar,
		&evento.EventoEstado,
		&evento.CantMeGusta,
		&evento.CantNoInteresa,
		&evento.CantVendidoTotal,
		&evento.TotalRecaudado,
		&evento.Estado,
		&evento.UsuarioCreacion,
		&evento.FechaCreacion,
		&evento.UsuarioModificacion,
		&evento.FechaModificacion,
	)

	if err != nil {
		return nil, err
	}

	return &evento, nil
}

func (e EventoModel) Patch(evento *Evento) error {
	query := `
		UPDATE evento
		SET organizador_id = $1, categoria_id = $2, titulo = $3, descripcion = $4,
			lugar = $5, evento_estado = $6, cant_me_gusta = $7, cant_no_interesa = $8,
			cant_vendido_total = $9, total_recaudado = $10, usuario_modificacion = $11,
			fecha_modificacion = NOW()
		WHERE evento_id = $12 AND estado = 1`

	args := []any{
		evento.OrganizadorID,
		evento.CategoriaID,
		evento.Titulo,
		evento.Descripcion,
		evento.Lugar,
		evento.EventoEstado,
		evento.CantMeGusta,
		evento.CantNoInteresa,
		evento.CantVendidoTotal,
		evento.TotalRecaudado,
		evento.UsuarioModificacion,
		evento.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := e.DB.QueryRowContext(ctx, query, args...).Scan(nil)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (e EventoModel) Delete(id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrRecordNotFound
	}

	// Soft delete - solo cambiamos el estado a 0
	query := `
		UPDATE evento
		SET estado = 0, fecha_modificacion = NOW()
		WHERE evento_id = $1 AND estado = 1`

	result, err := e.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (e EventoModel) GetAll(titulo string, filters []util.Filters) ([]*Evento, error) {
	query := `
		SELECT evento_id, organizador_id, categoria_id, titulo, descripcion, lugar,
			evento_estado, cant_me_gusta, cant_no_interesa, cant_vendido_total,
			total_recaudado, estado, usuario_creacion, fecha_creacion,
			usuario_modificacion, fecha_modificacion
		FROM evento
		WHERE (titulo ILIKE '%' || $1 || '%' OR $1 = '')
			AND estado = 1
		ORDER BY ` + filters[0].Sort + `
		LIMIT $2 OFFSET $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		titulo,
		filters[0].PageSize,
		(filters[0].Page - 1) * filters[0].PageSize,
	}

	var eventos []*Evento
	err := e.DB.SelectContext(ctx, &eventos, query, args...)

	if err != nil {
		return nil, err
	}

	return eventos, nil
}
