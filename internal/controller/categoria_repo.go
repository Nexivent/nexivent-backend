package data

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Loui27/nexivent-backend/internal/domain"
)

var _ domain.CategoriaRepository = (*CategoriaRepo)(nil)

type CategoriaRepo struct{ DB *sql.DB }

func NewCategoriaRepo(db *sql.DB) *CategoriaRepo { return &CategoriaRepo{DB: db} }

func (r *CategoriaRepo) Save(cont context.Context, c *domain.Categoria) error {
	if c.ID == 0 {
		cols := []string{"nombre", "descripcion"}
		vals := []any{c.Nombre, c.Descripcion}
		return InsertReturningID(cont, r.DB, "categorias", cols, vals, "id_categoria", &c.ID)
	}
	cols := []string{"nombre", "descripcion"}
	vals := []any{c.Nombre, c.Descripcion}
	return UpdateByID(cont, r.DB, "categorias", cols, vals, "id_categoria", c.ID)
}

func (r *CategoriaRepo) GetById(cont context.Context, id int) (*domain.Categoria, error) {
	const q = "Select id_categoria, nombre, descripcion FROM categorias WHERE id_categoria=$1"
	var out domain.Categoria
	if err := r.DB.QueryRowContext(cont, q, id).Scan(&out.ID, &out.Nombre, &out.Descripcion); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (r *CategoriaRepo) Delete(cont context.Context, id int) error {
	_, err := r.DB.ExecContext(cont, "DELETE FROM categorias WHERE id_categoria=$1", id)
	return err
}
