// evento.go
package model

import "time"

// EventDateView se utiliza para exponer fechas formateadas en respuestas JSON sin afectar el modelo persistente.
type EventDateView struct {
	IdFechaEvento int64  `json:"idFechaEvento"`
	IdFecha       int64  `json:"idFecha"`
	Fecha         string `json:"fecha"`
	HoraInicio    string `json:"horaInicio"`
	HoraFin       string `json:"horaFin"`
}

type Evento struct {
	ID                  int64 `gorm:"column:evento_id;primaryKey;autoIncrement"`
	OrganizadorID       int64
	CategoriaID         int64
	Titulo              string
	Descripcion         string
	Lugar               string
	EventoEstado        int16 `gorm:"default:0"`
	CantMeGusta         int   `gorm:"default:0"`
	CantNoInteresa      int   `gorm:"default:0"`
	CantVendidoTotal    int   `gorm:"default:0"`
	ImagenDescripcion   string
	ImagenPortada       string
	VideoPresentacion   string
	ImagenEscenario     string
	TotalRecaudado      float64
	Estado              int16 `gorm:"default:1"`
	UsuarioCreacion     *int64
	FechaCreacion       time.Time `gorm:"default:now()"`
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	Organizador *Usuario   `gorm:"foreignKey:OrganizadorID;references:ID"`
	Categoria   *Categoria `gorm:"foreignKey:CategoriaID;references:ID"`

	Comentarios []Comentario
	Sectores    []Sector          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TiposTicket []TipoDeTicket    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Perfiles    []PerfilDePersona `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Fechas      []EventoFecha     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// RELACIÓN 1–N: un evento tiene muchos cupones
	Cupones []Cupon `gorm:"foreignKey:EventoID;references:ID"`

	// EventDates expone fechas formateadas para respuestas sin tocar la estructura persistente.
	EventDates []EventDateView `gorm:"-" json:"eventDates"`
}

func (Evento) TableName() string { return "evento" }
