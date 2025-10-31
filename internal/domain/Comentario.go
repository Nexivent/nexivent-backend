package domain

type Comentario struct {
	IDComentario        int    `db:"id_comentario" json:"idComentario"`
	UsuarioComentadorID int    `db:"id_usuario_comentador" json:"usuarioComentadorId"`
	EventoComentadoID   int    `db:"id_evento_comentado" json:"eventoComentadoId"`
	Descripcion         string `db:"descripcion" json:"descripcion"`
}
