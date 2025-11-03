package domain

import "database/sql"

type Usuario struct {
	ID                    int64          `db:"usuario_id" json:"usuarioId"`
	Nombre                string         `db:"nombre" json:"nombre"`
	TipoDocumento         TipoDocumento  `db:"tipo_documento" json:"tipoDocumento"`
	NumDocumento          string         `db:"num_documento" json:"numDocumento"`
	Correo                string         `db:"correo" json:"correo"`
	Contrasenha           string         `db:"contrasenha" json:"-"`
	Telefono              sql.NullString `db:"telefono" json:"telefono,omitempty"`
	EstadoDeCuenta        int16          `db:"estado_de_cuenta" json:"estadoDeCuenta"`
	CodigoVerificacion    sql.NullString `db:"codigo_verificacion" json:"codigoVerificacion,omitempty"`
	FechaExpiracionCodigo sql.NullTime   `db:"fecha_expiracion_codigo" json:"fechaExpiracionCodigo,omitempty"`
}
