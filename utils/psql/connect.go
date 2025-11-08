	package psql

	import (
		"fmt"
		"log"
		"os"
		"time"

		"gorm.io/driver/postgres"
		"gorm.io/gorm"
		"gorm.io/gorm/logger"
	)

	// Create postgresql connection
	func CreateConnection(
		host string,
		user string,
		password string,
		dbname string,
		port string,
		enableLogs bool,
	) (*gorm.DB, error) {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s",
			host, user, password, dbname, port,
		)
		gormConfig := &gorm.Config{}
		if enableLogs {
			gormConfig.Logger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			})
		}

		return gorm.Open(postgres.Open(dsn), gormConfig)
	}
