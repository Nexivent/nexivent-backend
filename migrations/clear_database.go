package main

import (
	config "github.com/Loui27/nexivent-backend/internal/config"
	dao "github.com/Loui27/nexivent-backend/internal/dao/repository"
	"github.com/Loui27/nexivent-backend/logging"
	setupDB "github.com/Loui27/nexivent-backend/utils"
)

func main() {
	testLogger := logging.NewLoggerMock()
	envSettings := config.NuevoConfigEnv(testLogger)
	_, nexiventPsqlDB := dao.NewNexiventPsqlEntidades(testLogger, envSettings)

	setupDB.ClearPostgresqlDatabase(testLogger, nexiventPsqlDB, envSettings, nil)
}
