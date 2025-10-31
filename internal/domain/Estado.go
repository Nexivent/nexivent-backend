package domain

type Estado string

const (
	EstadoBorrador  Estado = "Borrador"
	EstadoPublicado Estado = "Publicado"
	EstadoCancelado Estado = "Cancelado"
)
