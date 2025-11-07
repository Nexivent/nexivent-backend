package repository

import (
	"fmt"

	config "github.com/Loui27/nexivent-backend/internal/config"
	model "github.com/Loui27/nexivent-backend/internal/dao/model"
	"github.com/Loui27/nexivent-backend/logging"
	psql "github.com/Loui27/nexivent-backend/utils/psql"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/gorm"
)

type NexiventPsqlEntidades struct {
	Logger logging.Logger
	Cupon  *Cupon
}

// Clase que crea colección de entidades para Nexivent Postgresql
func NewNexiventPsqlEntidades(
	logger logging.Logger,
	configEnv *config.ConfigEnv,
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
		Logger: logger,
		Cupon:  NewCuponController(logger, postgresqlDB),
	}, postgresqlDB
}

// Crear tablas dentro de la BD local
func crearTablas(astroCatPsqlDB *gorm.DB) {
	fmt.Println("Empezando creación de tablas...")

	// Crear tabla Cupon
	fmt.Println("Creando tabla Cupon...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Cupon{}); err != nil {
		fmt.Printf("Error creando tabla Cupon: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla Cupon creada exitosamente.")

	// Crear tabla Usuario
	fmt.Println("Creando tabla Usuario...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Usuario{}); err != nil {
		fmt.Printf("Error creando tabla Usuario: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla Usuario creada exitosamente.")

	// Crear tabla Categoria
	fmt.Println("Creando tabla Categoria...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Categoria{}); err != nil {
		fmt.Printf("Error creando tabla Categoria: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla Categoria creada exitosamente.")

	// Crear tabla Evento
	fmt.Println("Creando tabla Evento...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Evento{}); err != nil {
		fmt.Printf("Error creando tabla Evento: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla Evento creada exitosamente.")

	// Crear tabla Comentario
	fmt.Println("Creando tabla Comentario...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Comentario{}); err != nil {
		fmt.Printf("Error creando tabla Comentario: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla Comentario creada exitosamente.")

	// Crear tabla Sector
	fmt.Println("Creando tabla Sector...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Sector{}); err != nil {
		fmt.Printf("Error creando tabla Sector: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla Sector creada exitosamente.")

	// Crear tabla TipoDeTicket
	fmt.Println("Creando tabla TipoDeTicket...")
	if err := astroCatPsqlDB.AutoMigrate(&model.TipoDeTicket{}); err != nil {
		fmt.Printf("Error creando tabla TipoDeTicket: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla TipoDeTicket creada exitosamente.")

	// Crear tabla PerfilDePersona
	fmt.Println("Creando tabla PerfilDePersona...")
	if err := astroCatPsqlDB.AutoMigrate(&model.PerfilDePersona{}); err != nil {
		fmt.Printf("Error creando tabla PerfilDePersona: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla PerfilDePersona creada exitosamente.")

	// Crear tabla Tarifa
	fmt.Println("Creando tabla Tarifa...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Tarifa{}); err != nil {
		fmt.Printf("Error creando tabla Tarifa: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla Tarifa creada exitosamente.")

	// Crear tabla MetodoDePago
	fmt.Println("Creando tabla MetodoDePago...")
	if err := astroCatPsqlDB.AutoMigrate(&model.MetodoDePago{}); err != nil {
		fmt.Printf("Error creando tabla MetodoDePago: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla MetodoDePago creada exitosamente.")

	// Crear tabla OrdenDeCompra
	fmt.Println("Creando tabla OrdenDeCompra...")
	if err := astroCatPsqlDB.AutoMigrate(&model.OrdenDeCompra{}); err != nil {
		fmt.Printf("Error creando tabla OrdenDeCompra: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla OrdenDeCompra creada exitosamente.")

	// Crear tabla Fecha
	fmt.Println("Creando tabla Fecha...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Fecha{}); err != nil {
		fmt.Printf("Error creando tabla Fecha: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla Fecha creada exitosamente.")

	// Crear tabla EventoFecha
	fmt.Println("Creando tabla EventoFecha...")
	if err := astroCatPsqlDB.AutoMigrate(&model.EventoFecha{}); err != nil {
		fmt.Printf("Error creando tabla EventoFecha: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla EventoFecha creada exitosamente.")

	// Crear tabla Ticket
	fmt.Println("Creando tabla Ticket...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Ticket{}); err != nil {
		fmt.Printf("Error creando tabla Ticket: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla Ticket creada exitosamente.")

	// Crear tabla Cupon (otra vez por si la necesitas en otro contexto)
	fmt.Println("Creando tabla Cupon...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Cupon{}); err != nil {
		fmt.Printf("Error creando tabla Cupon: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla Cupon creada exitosamente.")

	// Crear tabla EventoCupon
	fmt.Println("Creando tabla EventoCupon...")
	if err := astroCatPsqlDB.AutoMigrate(&model.EventoCupon{}); err != nil {
		fmt.Printf("Error creando tabla EventoCupon: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla EventoCupon creada exitosamente.")

	// Crear tabla UsuarioCupon
	fmt.Println("Creando tabla UsuarioCupon...")
	if err := astroCatPsqlDB.AutoMigrate(&model.UsuarioCupon{}); err != nil {
		fmt.Printf("Error creando tabla UsuarioCupon: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla UsuarioCupon creada exitosamente.")

	// Crear tabla ComprobanteDePago
	fmt.Println("Creando tabla ComprobanteDePago...")
	if err := astroCatPsqlDB.AutoMigrate(&model.ComprobanteDePago{}); err != nil {
		fmt.Printf("Error creando tabla ComprobanteDePago: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla ComprobanteDePago creada exitosamente.")

	// Crear tabla Rol
	fmt.Println("Creando tabla Rol...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Rol{}); err != nil {
		fmt.Printf("Error creando tabla Rol: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla Rol creada exitosamente.")

	// Crear tabla RolUsuario
	fmt.Println("Creando tabla RolUsuario...")
	if err := astroCatPsqlDB.AutoMigrate(&model.RolUsuario{}); err != nil {
		fmt.Printf("Error creando tabla RolUsuario: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla RolUsuario creada exitosamente.")

	// Crear tabla Notificacion
	fmt.Println("Creando tabla Notificacion...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Notificacion{}); err != nil {
		fmt.Printf("Error creando tabla Notificacion: %v\n", err)
		panic(err)
	}
	fmt.Println("Tabla Notificacion creada exitosamente.")

	fmt.Println("Todas las tablas fiueron creadas exitosamente.")
}

// Reiniciar la base de datos local borrando todas las tablas
func BorrarTodasLasTablas(nexiventPsqlDB *gorm.DB) {
	// Disable foreign keys
	nexiventPsqlDB.Exec("SET CONSTRAINTS ALL DEFERRED")

	// Drop all tables in reverse order of dependencies
	tablesToDrop := []string{
		"rol_usuario",
		"usuario_cupon",
		"evento_cupon",
		"ticket",
		"comprobante_de_pago",
		"evento_fecha",
		"fecha",
		"tarifa",
		"sector",
		"tipo_de_ticket",
		"perfil_de_persona",
		"comentario",
		"orden_de_compra",
		"metodo_de_pago",
		"evento",
		"cupon",
		"rol",
		"notificacion",
		"categoria",
		"usuario",
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
