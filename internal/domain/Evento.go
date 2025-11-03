package domain

import "github.com/shopspring/decimal"

type Evento struct {
	ID               int64           `db:"evento_id" json:"eventoId"`
	OrganizadorID    Usuario         `db:"organizador_id" json:"organizadorId"`
	CategoriaID      Categoria       `db:"categoria_id" json:"categoriaId"`
	Titulo           string          `db:"titulo" json:"titulo"`
	Descripcion      string          `db:"descripcion" json:"descripcion"`
	Lugar            string          `db:"lugar" json:"lugar"`
	EventoEstado     int16           `db:"evento_estado" json:"eventoEstado"`
	CantMeGusta      int             `db:"cant_me_gusta" json:"cantMeGusta"`
	CantNoInteresa   int             `db:"cant_no_interesa" json:"cantNoInteresa"`
	CantVendidoTotal int             `db:"cant_vendido_total" json:"cantVendidoTotal"`
	TotalRecaudado   decimal.Decimal `db:"total_recaudado" json:"totalRecaudado"` // o int64 en centavos
	Activo           int16           `db:"activo" json:"activo"`
}
