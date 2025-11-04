package main

import (
	"github.com/Loui27/nexivent-backend/logging"
	domain "github.com/Loui27/nexivent-backend/internal/domain"
	dao "github.com/Loui27/nexivent-backend/internal/dao/nexivent-psql"
	setupDB "github.com/Loui27/nexivent-backend/utils"
)

func main() {
	testLogger := logging.NewLoggerMock()
	envSettings := domain.NuevoConfigEnv(testLogger)
	_, nexiventPsqlDB := dao.NewNexiventPsqlEntidades(testLogger, envSettings)

	setupDB.ClearPostgresqlDatabase(testLogger, nexiventPsqlDB, envSettings, nil)
}