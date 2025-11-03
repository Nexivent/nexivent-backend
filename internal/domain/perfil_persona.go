package domain

type PerfilDePersona struct {
	ID     int64  `db:"perfil_de_persona_id" json:"perfilDePersonaId"`
	Evento Evento `db:"-" json:"evento"`
	Nombre string `db:"nombre" json:"nombre"`
	Activo int16  `db:"activo" json:"activo"`
}
