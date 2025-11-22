package api

import (
	"github.com/Nexivent/nexivent-backend/internal/application/controller"
	config "github.com/Nexivent/nexivent-backend/internal/config"
	"github.com/Nexivent/nexivent-backend/logging"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Api struct {
	Logger        logging.Logger
	BllController *controller.ControllerCollection
	ConfigEnv     *config.ConfigEnv
	Echo          *echo.Echo
}

/*
Creates a new Api server with
- Logger provided by input
- BllController as new bll controller collection
- ConfigEnv as new config env settings
*/
func NewApi(
	logger logging.Logger,
	configEnv *config.ConfigEnv,
) (*Api, *gorm.DB) {
	bllController, nexiventPsqlDB := controller.NewControllerCollection(logger, configEnv)

	return &Api{
		Logger:        logger,
		BllController: bllController,
		ConfigEnv:     configEnv,
		Echo:          echo.New(),
	}, nexiventPsqlDB
}

// @title Nexivent API
// @version 1.0
// @description Nexivent Event Management API
// @BasePath /
func RunService(configEnv *config.ConfigEnv, logger logging.Logger) {
	api, _ := NewApi(logger, configEnv)
	api.RunApi(configEnv)
}
