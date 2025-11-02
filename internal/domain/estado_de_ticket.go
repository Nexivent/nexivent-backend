package domain

type EstadoDeTicket string

const (
	Disponible EstadoDeTicket = "Disponible"
	Vendido    EstadoDeTicket = "Vendido"
	Usado      EstadoDeTicket = "Usado"
	Cancelado  EstadoDeTicket = "Cancelado"
)
