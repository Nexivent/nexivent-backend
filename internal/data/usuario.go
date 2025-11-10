package data

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
	"github.com/google/uuid"
)

type Usuario struct {
	ID                    uuid.UUID
	Nombre                string
	TipoDocumento         string
	NumDocumento          string
	Correo                string
	Contrasenha           string
	Telefono              *string
	EstadoDeCuenta        util.Estado
	CodigoVerificacion    *string
	FechaExpiracionCodigo *time.Time
	UsuarioCreacion       *uuid.UUID
	FechaCreacion         time.Time
	UsuarioModificacion   *uuid.UUID
	FechaModificacion     *time.Time
	Estado                util.Estado

	Comentarios    []Comentario
	Ordenes        []OrdenDeCompra
	RolesAsignados []RolUsuario
	Cupones        []UsuarioCupon
}

func (Usuario) TableName() string { return "usuario" }
