package utils

import (
	"context"
	"testing"
	"time"

	config "github.com/Nexivent/nexivent-backend/internal/config"
	model "github.com/Nexivent/nexivent-backend/internal/dao/model"
	"github.com/Nexivent/nexivent-backend/logging"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CustomLogger struct{}

// Helper function to create string pointers
func strPtr(s string) *string {
	return &s
}

func (l *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	// Ignore info logs
}

func (l *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	// Ignore warning logs
}

func (l *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	// Ignore error logs
}

func (l *CustomLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (string, int64),
	err error,
) {
	// Ignore trace logs
}

func ClearPostgresqlDatabaseTesting(
	appLogger logging.Logger,
	nexiventDB *gorm.DB,
	envSetting *config.ConfigEnv,
	t *testing.T,
) {
	if envSetting.PostgresHost != "localhost" {
		msg := "Not allow clear Levels Postgres DB into instance different to localhost"
		if t == nil {
			appLogger.Panicf(
				"%s. This function should only be used for tests in local environment",
				msg,
			)
		} else {
			t.Fatalf("%s. This function should only be used for tests in local environment", msg)
		}
		return
	}

	if nexiventDB != nil {
		// fmt.Println("...Clearing AstroCatPsql database (hard delete)...")

		originalLogger := nexiventDB.Logger
		if !envSetting.EnableSqlLogs {
			nexiventDB.Logger = originalLogger.LogMode(logger.Silent)
		}

		// Start a transaction
		tx := nexiventDB.Begin()

		// Disable foreign key constraints temporarily
		tx.Exec("SET CONSTRAINTS ALL DEFERRED")

		// First delete tables that have references to other tables
		tablesToClear := []struct {
			name  string
			model any
		}{
			// First delete tables with foreign key dependencies
			{"rol_usuario", &model.RolUsuario{}},
			{"usuario_cupon", &model.UsuarioCupon{}},
			//{"evento_cupon", &model.EventoCupon{}},
			{"ticket", &model.Ticket{}},
			{"comprobante_de_pago", &model.ComprobanteDePago{}},
			{"evento_fecha", &model.EventoFecha{}},
			{"fecha", &model.Fecha{}},
			{"tarifa", &model.Tarifa{}},
			{"sector", &model.Sector{}},
			{"tipo_de_ticket", &model.TipoDeTicket{}},
			{"perfil_de_persona", &model.PerfilDePersona{}},
			{"comentario", &model.Comentario{}},
			{"orden_de_compra", &model.OrdenDeCompra{}},
			{"metodo_de_pago", &model.MetodoDePago{}},
			{"evento", &model.Evento{}},
			{"cupon", &model.Cupon{}},
			{"rol", &model.Rol{}},
			{"notificacion", &model.Notificacion{}},
			{"categoria", &model.Categoria{}},
			{"usuario", &model.Usuario{}}, // Clear audit logs first to avoid FK constraints
		}

		for _, table := range tablesToClear {
			appLogger.Infof("Attempting to hard delete all records from %s table...", table.name)
			if err := tx.Unscoped().Where("true").Delete(table.model).Error; err != nil {
				tx.Rollback()
				appLogger.Errorf("Error clearing %s table: %v", table.name, err)
				return
			}
		}

		// Reactivate foreign key constraints
		tx.Exec("SET CONSTRAINTS ALL IMMEDIATE")

		// Confirm the transaction
		if err := tx.Commit().Error; err != nil {
			appLogger.Errorf("Error committing transaction: %v", err)
			return
		}

		if !envSetting.EnableSqlLogs {
			nexiventDB.Logger = originalLogger
		}

	} else {
		appLogger.Warn("nexiventDB is nil, skipping database clearing.")
	}
}

// Remove all data from AstroCatPsql db.
//   - Note: Only use for tests
func ClearPostgresqlDatabase(
	appLogger logging.Logger,
	nexiventDB *gorm.DB,
	envSetting *config.ConfigEnv,
	t *testing.T,
) {
	if envSetting.PostgresHost != "localhost" {
		msg := "Not allow clear Levels Postgres DB into instance different to localhost"
		if t == nil {
			appLogger.Panicf(
				"%s. This function should only be used for tests in local environment",
				msg,
			)
		} else {
			t.Fatalf("%s. This function should only be used for tests in local environment", msg)
		}
		return
	}

	if nexiventDB != nil {
		// fmt.Println("...Clearing AstroCatPsql database (hard delete)...")

		originalLogger := nexiventDB.Logger
		if !envSetting.EnableSqlLogs {
			nexiventDB.Logger = originalLogger.LogMode(logger.Silent)
		}

		// Start a transaction
		tx := nexiventDB.Begin()

		// Disable foreign key constraints temporarily
		tx.Exec("SET CONSTRAINTS ALL DEFERRED")

		// First delete tables that have references to other tables
		tablesToClear := []struct {
			name  string
			model any
		}{
			// First delete tables with foreign key dependencies
			{"rol_usuario", &model.RolUsuario{}},
			{"usuario_cupon", &model.UsuarioCupon{}},
			//{"evento_cupon", &model.EventoCupon{}},
			{"ticket", &model.Ticket{}},
			{"comprobante_de_pago", &model.ComprobanteDePago{}},
			{"evento_fecha", &model.EventoFecha{}},
			{"fecha", &model.Fecha{}},
			{"tarifa", &model.Tarifa{}},
			{"sector", &model.Sector{}},
			{"tipo_de_ticket", &model.TipoDeTicket{}},
			{"perfil_de_persona", &model.PerfilDePersona{}},
			{"comentario", &model.Comentario{}},
			{"orden_de_compra", &model.OrdenDeCompra{}},
			{"metodo_de_pago", &model.MetodoDePago{}},
			{"evento", &model.Evento{}},
			{"cupon", &model.Cupon{}},
			{"rol", &model.Rol{}},
			{"notificacion", &model.Notificacion{}},
			{"categoria", &model.Categoria{}},
			{"usuario", &model.Usuario{}},
		}

		for _, table := range tablesToClear {
			appLogger.Infof("Attempting to hard delete all records from %s table...", table.name)
			if err := tx.Unscoped().Where("true").Delete(table.model).Error; err != nil {
				tx.Rollback()
				appLogger.Errorf("Error clearing %s table: %v", table.name, err)
				return
			}
		}

		// Reactivate foreign key constraints
		tx.Exec("SET CONSTRAINTS ALL IMMEDIATE")

		// Confirm the transaction
		if err := tx.Commit().Error; err != nil {
			appLogger.Errorf("Error committing transaction: %v", err)
			return
		}

		if !envSetting.EnableSqlLogs {
			nexiventDB.Logger = originalLogger
		}
	} else {
		appLogger.Warn("nexiventDB is nil, skipping database clearing.")
	}
}
