package schemas

type RolUsuarioResponse struct {
	IDUsuario int64 `json:"idUsuario"`
	Roles     []RolResponse `json:"roles"` 
}


type RolUsuarioRequest struct {
	IDUsuario int64 `json:"idUsuario"`
	IDRol     int64 `json:"idRol"`
}

type UsuarioRolResponse struct {
	IDUsuario  int64    `json:"idUsuario"`
	Nombre     string   `json:"nombre"`
	Correo     string   `json:"correo"`
	Estado     int16    `json:"estado"`
}
