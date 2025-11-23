package schemas

type RolUsuarioResponse struct {
	IDUsuario int64 `json:"idUsuario"`
	Roles     []RolResponse `json:"roles"` 
}


type RolUsuarioRequest struct {
	IDUsuario int64 `json:"idUsuario"`
	IDRol     int64 `json:"idRol"`
}
