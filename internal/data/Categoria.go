package data

type Comentario struct {
	IDComentario      int     `db:"id_comentario" json:"idComentario"`
	UsuarioComentador Usuario `db:"-" json:"usuarioComentador"`
	EventoComentado   Evento  `db:"-" json:"eventoComentado"`
	Descripcion       string  `db:"descripcion" json:"descripcion"`
}

