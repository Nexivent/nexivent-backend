package data

import "github.com/google/uuid"

type UsuarioCupon struct {
	CuponID   uuid.UUID
	UsuarioID uuid.UUID
	CantUsada int64

	Cupon   *Cupon
	Usuario *Usuario
}


