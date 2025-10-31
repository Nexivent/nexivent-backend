package data

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Loui27/nexivent-backend/internal/domain"
)

type CategoriaRepo struct{ DB *sql.DB }

func NewCategoriaRepo(db *sql.DB) *CategoriaRepo { return &CategoriaRepo{DB: db} }

func (r *CategoriaRepo) Save(cont context.Context, c *domain.Categoria) error {
	if c.IDCategoria == 0 {
		cols := []string{"nombre", "descripcion"}
		vals := []any{c.Nombre, c.Descripcion}
		return InsertReturningID(cont, r.DB, "categorias", cols, vals, "id_categoria", &c.IDCategoria)
	}
	cols := []string{"nombre", "descripcion"}
	vals := []any{c.Nombre, c.Descripcion}
	return UpdateByID(cont, r.DB, "categorias", cols, vals, "id_categoria", c.IDCategoria)
}

func (r *CategoriaRepo) GetById(cont context.Context, id int64) (*domain.Categoria, error) {
	const q = "Select id_categoria, nombre, descripcion FROM categorias WHERE id_categoria=$1"
	var out domain.Categoria
	if err := r.DB.QueryRowContext(cont, q, id).Scan(&out.IDCategoria, &out.Nombre, &out.Descripcion); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (r *CategoriaRepo) Delete(cont context.Context, id int64) error {
	_, err := r.DB.ExecContext(cont, "DELETE FROM categorias WHERE id_categoria=$1", id)
	return err
}
