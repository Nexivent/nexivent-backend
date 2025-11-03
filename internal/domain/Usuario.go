package domain

import (
	"database/sql"
	"time"
)

type Usuario struct {
	ID                    int64          `db:"usuario_id" json:"id"`
	Nombre                string         `db:"nombre" json:"nombre"`
	TipoDocumento         TipoDocumento  `db:"tipo_documento" json:"tipoDocumento"`
	NumDocumento          string         `db:"num_documento" json:"numDocumento"`
	Correo                string         `db:"correo" json:"correo"`
	Contrasenha           string         `db:"contrasenha" json:"-"`
	Telefono              sql.NullString `db:"telefono" json:"telefono,omitempty"`
	EstadoDeCuenta        int16          `db:"estado_de_cuenta" json:"estadoDeCuenta"`
	CodigoVerificacion    sql.NullString `db:"codigo_verificacion" json:"codigoVerificacion,omitempty"`
	FechaExpiracionCodigo sql.NullTime   `db:"fecha_expiracion_codigo" json:"fechaExpiracionCodigo,omitempty"`
	UsuarioCreacionID     sql.NullInt64  `db:"usuario_creacion" json:"usuarioCreacionId,omitempty"`
	FechaCreacion         time.Time      `db:"fecha_creacion" json:"fechaCreacion"`
	UsuarioModificacionID sql.NullInt64  `db:"usuario_modificacion" json:"usuarioModificacionId,omitempty"`
	FechaModificacion     sql.NullTime   `db:"fecha_modificacion" json:"fechaModificacion,omitempty"`
	Activo                int16          `db:"activo" json:"activo"`
}
a