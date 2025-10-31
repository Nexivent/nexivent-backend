package data

type Usuario struct {
	IDUsuario   int    `db:"id_usuario" json:"idUsuario"`
	Nombre      string `db:"nombre" json:"nombre"`
	Contrasena  string `db:"contrasena" json:"-"`
	Correo      string `db:"correo" json:"correo"`
	NumTelefono string `db:"num_telefono" json:"numTelefono"`
	RolUsuario  Rol    `db:"rol_usuario" json:"rolUsuario"`
}
