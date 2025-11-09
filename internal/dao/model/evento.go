package model

import "time"

type Evento struct {
	ID                  int64 `gorm:"column:evento_id;primaryKey;autoIncrement"`
	OrganizadorID       int64
	CategoriaID         int64
	Titulo              string
	Descripcion         string
	Lugar               string
	EventoEstado        int16
	CantMeGusta         int
	CantNoInteresa      int
	CantVendidoTotal    int
	TotalRecaudado      float64
	Estado              int16
	UsuarioCreacion     int64
	FechaCreacion       time.Time
	UsuarioModificacion *int64
	FechaModificacion   *time.Time

	Organizador *Usuario   `gorm:"foreignKey:OrganizadorID"`
	Categoria   *Categoria `gorm:"foreignKey:CategoriaID"`

	Comentarios []Comentario
	Sectores    []Sector          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TiposTicket []TipoDeTicket    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Perfiles    []PerfilDePersona `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Fechas      []EventoFecha     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Evento) TableName() string { return "evento" }
