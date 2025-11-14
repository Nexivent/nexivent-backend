package model

import (
	"regexp"
	"time"

	"github.com/Nexivent/nexivent-backend/internal/data/model/util"
	"github.com/Nexivent/nexivent-backend/internal/validator"
)

type Usuario struct {
	ID                    uint64      `gorm:"column:usuario_id;primaryKey;autoIncrement" json:"id"`
	Nombre                string      `gorm:"column:nombre" json:"nombre"`
	TipoDocumento         string      `gorm:"column:tipo_documento;uniqueIndex:uq_usuario_doc" json:"tipoDocumento"`
	NumDocumento          string      `gorm:"column:num_documento;uniqueIndex:uq_usuario_doc" json:"numDocumento"`
	Correo                string      `gorm:"column:correo;uniqueIndex" json:"correo"`
	password              []byte      `gorm:"column:password" json:"-"`
	Telefono              *string     `gorm:"column:telefono" json:"telefono,omitempty"`
	EstadoDeCuenta        util.Estado `gorm:"column:estado_de_cuenta;default:0" json:"estadoDeCuenta"`
	CodigoVerificacion    *string     `gorm:"column:codigo_verificacion" json:"-"`
	FechaExpiracionCodigo *time.Time  `gorm:"column:fecha_expiracion_codigo" json:"-"`
	UsuarioCreacion       *uint64     `gorm:"column:usuario_creacion" json:"-"`
	FechaCreacion         time.Time   `gorm:"column:fecha_creacion;default:now()" json:"-"`
	UsuarioModificacion   *uint64     `gorm:"column:usuario_modificacion" json:"-"`
	FechaModificacion     *time.Time  `gorm:"column:fecha_modificacion" json:"-"`
	Estado                util.Estado `gorm:"column:estado;default:1" json:"-"`

	Comentarios    []Comentario    `json:"comentarios,"`
	Ordenes        []OrdenDeCompra `json:"ordenes"`
	RolesAsignados []RolUsuario    `json:"roles"`
	Cupones        []UsuarioCupon  `json:"cupones"`
}

func (Usuario) TableName() string { return "usuario" }

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateUsuario(v *validator.Validator, usuario *Usuario) {
	// Validar Nombre
	v.Check(usuario.Nombre != "", "nombre", "el nombre es obligatorio")
	v.Check(len(usuario.Nombre) <= 100, "nombre", "el nombre no debe exceder 100 caracteres")

	// Validar TipoDocumento
	v.Check(usuario.TipoDocumento != "", "tipoDocumento", "el tipo de documento es obligatorio")
	v.Check(usuario.TipoDocumento == "DNI" || usuario.TipoDocumento == "CE" || usuario.TipoDocumento == "RUC",
		"tipoDocumento", "el tipo de documento debe ser DNI, CE o RUC")

	// Validar NumDocumento
	v.Check(usuario.NumDocumento != "", "numDocumento", "el número de documento es obligatorio")
	v.Check(len(usuario.NumDocumento) <= 20, "numDocumento", "el número de documento no debe exceder 20 caracteres")

	// Validar Correo
	v.Check(usuario.Correo != "", "correo", "el correo es obligatorio")
	v.Check(emailRegex.MatchString(usuario.Correo), "correo", "el correo debe tener un formato válido")
	v.Check(len(usuario.Correo) <= 100, "correo", "el correo no debe exceder 100 caracteres")

	// Validar password
	v.Check(len(usuario.password) != 0, "contrasenha", "la contraseña debe tener al menos 6 caracteres")

	// Validar Telefono (si está presente)
	if usuario.Telefono != nil && *usuario.Telefono != "" {
		v.Check(len(*usuario.Telefono) <= 15, "telefono", "el teléfono no debe exceder 15 caracteres")
	}
}
