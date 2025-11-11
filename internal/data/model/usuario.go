package model

import (
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/util"
)

type Usuario struct {
	ID                    uint64
	Nombre                string
	TipoDocumento         string
	NumDocumento          string
	Correo                string
	Contrasenha           string
	Telefono              *string
	EstadoDeCuenta        util.Estado
	CodigoVerificacion    *string
	FechaExpiracionCodigo *time.Time
	UsuarioCreacion       *uint64
	FechaCreacion         time.Time
	UsuarioModificacion   *uint64
	FechaModificacion     *time.Time
	Estado                util.Estado

	Comentarios    []Comentario
	Ordenes        []OrdenDeCompra
	RolesAsignados []RolUsuario
	Cupones        []UsuarioCupon
}

func (Usuario) TableName() string { return "usuario" }
