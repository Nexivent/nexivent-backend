package domain

type Sector struct {
	ID            int64  `db:"sector_id" json:"sectorId"`
	Evento        Evento `db:"-" json:"eventoId"`
	SectorTipo    string `db:"sector_tipo" json:"sectorTipo"`
	TotalEntradas int    `db:"total_entradas" json:"totalEntradas"`
	CantVendidas  int    `db:"cant_vendidas" json:"cantVendidas"`
	Activo        int16  `db:"activo" json:"activo"`
}
