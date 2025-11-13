package repository

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Repository struct {
	DB              *gorm.DB
	Eventos         EventoSQL
	EventoFechas    EventoFecha
	Fechas          Fecha
	Categorias      Categoria
	Usuarios        Usuario
	Roles           Rol
	RolesUsuarios   RolUsuarioRepo
	Cupones         Cupon
	MetodosPago     MetodoDePago
	Tarifas         Tarifa
	OrdenesCompra   OrdenDeCompra
	Tickets         Ticket
	Comentarios     Comentario
	Sectores        Sector
	TiposDeTicket   TipoDeTicket
	PerfilesPersona PerfilDePersona
	Comprobantes    ComprobanteDePago
	Notificaciones  Notificacion
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{
		DB:              db,
		Eventos:         EventoSQL{DB: db},
		EventoFechas:    EventoFecha{DB: db},
		Fechas:          Fecha{DB: db},
		Categorias:      Categoria{DB: db},
		Usuarios:        Usuario{DB: db},
		Roles:           Rol{DB: db},
		RolesUsuarios:   RolUsuarioRepo{DB: db},
		Cupones:         Cupon{DB: db},
		MetodosPago:     MetodoDePago{DB: db},
		Tarifas:         Tarifa{DB: db},
		OrdenesCompra:   OrdenDeCompra{DB: db},
		Tickets:         Ticket{DB: db},
		Comentarios:     Comentario{DB: db},
		Sectores:        Sector{DB: db},
		TiposDeTicket:   TipoDeTicket{DB: db},
		PerfilesPersona: PerfilDePersona{DB: db},
		Comprobantes:    ComprobanteDePago{DB: db},
		Notificaciones:  Notificacion{DB: db},
	}
}

// WithTx crea una nueva instancia del Repository usando una transacción
func (r *Repository) WithTx(tx *gorm.DB) Repository {
	return Repository{
		DB:              tx,
		Eventos:         EventoSQL{DB: tx},
		EventoFechas:    EventoFecha{DB: tx},
		Fechas:          Fecha{DB: tx},
		Categorias:      Categoria{DB: tx},
		Usuarios:        Usuario{DB: tx},
		Roles:           Rol{DB: tx},
		RolesUsuarios:   RolUsuarioRepo{DB: tx},
		Cupones:         Cupon{DB: tx},
		MetodosPago:     MetodoDePago{DB: tx},
		Tarifas:         Tarifa{DB: tx},
		OrdenesCompra:   OrdenDeCompra{DB: tx},
		Tickets:         Ticket{DB: tx},
		Comentarios:     Comentario{DB: tx},
		Sectores:        Sector{DB: tx},
		TiposDeTicket:   TipoDeTicket{DB: tx},
		PerfilesPersona: PerfilDePersona{DB: tx},
		Comprobantes:    ComprobanteDePago{DB: tx},
		Notificaciones:  Notificacion{DB: tx},
	}
}
