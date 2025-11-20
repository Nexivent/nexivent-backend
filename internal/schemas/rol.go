package schemas

import (
	"time"
)

type RolRequest struct {
	//ID                  int64  `gorm:"column:rol_id;primaryKey;autoIncrement"`
	Nombre              string `json:"nombre"`
	UsuarioCreacion     *int64
	FechaCreacion       time.Time 
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	//Usuarios []RolUsuario
}

type RolResponse struct {
	ID                  int64  `json:"id"`
	Nombre              string `json:"nombre"`
	//UsuarioCreacion     *int64
	//FechaCreacion       time.Time `json:"fecha"`
	//UsuarioModificacion *int64
	//FechaModificacion   *time.Time

	//Usuarios []RolUsuario
}
