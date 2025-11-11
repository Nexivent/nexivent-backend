package repository

import (
	"time"

	"github.com/Loui27/nexivent-backend/internal/dao/model"
	"github.com/Loui27/nexivent-backend/logging"
	"gorm.io/gorm"
)

type Tarifa struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewTarifaController(
	logger logging.Logger,
	postgresqlDB *gorm.DB,
) *Tarifa {
	return &Tarifa{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

func (t *Tarifa) CrearTarifa(tarifa *model.Tarifa) error {
	res := t.PostgresqlDB.Create(tarifa)
	if res.Error != nil {
		t.logger.Errorf("CrearTarifa: %v", res.Error)
		return res.Error
	}
	return nil
}

func (t *Tarifa) ActualizarTarifa(tarifa *model.Tarifa) error {
	res := t.PostgresqlDB.Save(tarifa)
	if res.Error != nil {
		t.logger.Errorf("ActualizarTarifa: %v", res.Error)
		return res.Error
	}
	return nil
}

// ObtenerTarifasPorIDs: trae tarifas con estado=1 por sus IDs (sin validar evento/fecha).
func (t *Tarifa) ObtenerTarifasPorIDs(ids []int64) ([]*model.Tarifa, error) {
	if len(ids) == 0 {
		return []*model.Tarifa{}, nil
	}
	var list []*model.Tarifa
	res := t.PostgresqlDB.
		Where("tarifa_id IN ?", ids).
		Where("estado = 1").
		Find(&list)
	if res.Error != nil {
		t.logger.Errorf("ObtenerTarifasPorIDs: %v", res.Error)
		return nil, res.Error
	}
	return list, nil
}

// ObtenerTarifasValidasParaFechaEvento:
// - Filtra por ids de tarifa (si se proveen)
// - Valida que la tarifa esté activa (tarifa.estado = 1)
// - Valida que el tipo_de_ticket esté activo (tdt.estado = 1)
// - Valida que la tarifa pertenezca al evento (tdt.evento_id = idEvento)
// - Valida que la fecha (solo parte de día) esté entre [tdt.fecha_ini, tdt.fecha_fin]
// - Devuelve además los datos de Sector (para verificar stock en el BO)
func (t *Tarifa) ObtenerTarifasValidasParaFechaEvento(
	ids []int64,
	idEvento int64,
	fecha time.Time,
) ([]*model.Tarifa, error) {

	fechaDia := fecha.Format("2006-01-02") // comparar solo día contra DATE

	q := t.PostgresqlDB.
		Table("tarifa AS tf").
		Select(`
			tf.tarifa_id, tf.sector_id, tf.tipo_de_ticket_id, tf.perfil_de_persona_id,
			tf.precio, tf.estado,
			tf.usuario_creacion, tf.fecha_creacion, tf.usuario_modificacion, tf.fecha_modificacion
		`).
		Joins("JOIN tipo_de_ticket AS tdt ON tdt.tipo_de_ticket_id = tf.tipo_de_ticket_id").
		Where("tf.estado = 1").
		Where("tdt.estado = 1").
		Where("tdt.evento_id = ?", idEvento).
		Where("?::date BETWEEN tdt.fecha_ini AND tdt.fecha_fin", fechaDia)

	if len(ids) > 0 {
		q = q.Where("tf.tarifa_id IN ?", ids)
	}

	var list []*model.Tarifa
	if err := q.Find(&list).Error; err != nil {
		t.logger.Errorf("ObtenerTarifasValidasParaFechaEvento(idEvento=%d, fecha=%s): %v", idEvento, fechaDia, err)
		return nil, err
	}

	return list, nil
}

// MapTarifaPrecioSector:
// Construye mapas útiles para el BO:
//   - precioPorTarifa[tarifaID] = precio
//   - sectorPorTarifa[tarifaID] = sectorID
func (t *Tarifa) MapTarifaPrecioSector(ids []int64) (map[int64]float64, map[int64]int64, error) {
	outPrecio := make(map[int64]float64, len(ids))
	outSector := make(map[int64]int64, len(ids))

	if len(ids) == 0 {
		return outPrecio, outSector, nil
	}

	var rows []struct {
		ID       int64   `gorm:"column:tarifa_id"`
		Precio   float64 `gorm:"column:precio"`
		SectorID int64   `gorm:"column:sector_id"`
	}
	res := t.PostgresqlDB.
		Table("tarifa").
		Select("tarifa_id, precio, sector_id").
		Where("tarifa_id IN ?", ids).
		Where("estado = 1").
		Find(&rows)
	if res.Error != nil {
		t.logger.Errorf("MapTarifaPrecioSector: %v", res.Error)
		return nil, nil, res.Error
	}

	for _, r := range rows {
		outPrecio[r.ID] = r.Precio
		outSector[r.ID] = r.SectorID
	}
	return outPrecio, outSector, nil
}
