package nexiventpsql

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/Loui27/nexivent-backend/logging"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	domain "github.com/Loui27/nexivent-backend/internal/domain"
	model "github.com/Loui27/nexivent-backend/internal/dao/model"
	psql "github.com/Loui27/nexivent-backend/utils/psql"
)

type NexiventPsqlEntidades struct {
	Logger               logging.Logger
	Cupon 							*Cupon
}

// Clase que crea colección de entidades para Nexivent Postgresql
func NewNexiventPsqlEntidades(
	logger logging.Logger,
	configEnv *domain.ConfigEnv,
) (*NexiventPsqlEntidades, *gorm.DB) {
	postgresqlDB, err := psql.CreateConnection(
		configEnv.PostgresHost,
		configEnv.PostgresUser,
		configEnv.PostgresPassword,
		configEnv.PostgresDBName,
		configEnv.PostgresPort,
		configEnv.EnableSqlLogs,
	)
	if err != nil {
		logger.Panicln("Failed to connect to AstroCat Postgresql database")
	}

	if err := postgresqlDB.Use(otelgorm.NewPlugin()); err != nil {
		logger.Panicln("Failed to instrument AstroCat Postgresql database")
	}

	crearTablas(postgresqlDB)

	return &NexiventPsqlEntidades{
		Logger:               logger,
		Cupon:            		NewCuponController(logger, postgresqlDB),
	}, postgresqlDB
}

// Crear tablas dentro de la BD local
func crearTablas(astroCatPsqlDB *gorm.DB) {
	fmt.Println("Empezando creación de tablas...")

	fmt.Println("Creando tabla Cupon...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Cupon{}); err != nil {
		fmt.Printf("Error creando tabla Cupon: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla Cupon creada exitosamente.")
	// Crear nueva tabla
	// fmt.Println("Creando tabla Cupon...")
	// if err := astroCatPsqlDB.AutoMigrate(&model.Cupon{}); err != nil {
	// 	fmt.Printf("Error creando tabla Cupon: %v\n", err)
	// 	panic(err)
	// }
	// fmt.Println("Tabla Cupon creada exitosamente.")

	fmt.Println("Todas las tablas fiueron creadas exitosamente.")
}

// Reiniciar la base de datos local borrando todas las tablas
func BorrarTodasLasTablas(nexiventPsqlDB *gorm.DB) {
	// Disable foreign keys
	nexiventPsqlDB.Exec("SET CONSTRAINTS ALL DEFERRED")

	// Drop all tables in reverse order of dependencies
	tablesToDrop := []string{
		"nexivent_cupon",
	}

	for _, tableName := range tablesToDrop {
		if err := nexiventPsqlDB.Exec("DROP TABLE IF EXISTS " + tableName + " CASCADE").Error; err != nil {
			fmt.Printf("Warning: Error dropping table %s: %v\n", tableName, err)
		} else {
			fmt.Printf("Dropped table: %s\n", tableName)
		}
	}

	// Reactivate foreign key constraints
	nexiventPsqlDB.Exec("SET CONSTRAINTS ALL IMMEDIATE")
}
