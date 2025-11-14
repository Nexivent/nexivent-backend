package api

import (
	"github.com/Loui27/nexivent-backend/logging"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	//"onichankimochi.com/astro_cat_backend/src/server/bll/controller"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Api struct {
	Logger        logging.Logger
	BllController *controller.ControllerCollection
	EnvSettings   *schemas.EnvSettings
	Echo          *echo.Echo
}

/*
Creates a new Api server with
- Logger provided by input
- BllController as new bll controller collection
- EnvSettings as new env settings provided by .env file
*/
func NewApi(
	logger logging.Logger,
	envSettings *schemas.EnvSettings,
) (*Api, *gorm.DB) {
	bllController, astroCatPsqlDB := controller.NewControllerCollection(logger, envSettings)

	return &Api{
		Logger:        logger,
		BllController: bllController,
		EnvSettings:   envSettings,
		Echo:          echo.New(),
	}, astroCatPsqlDB
}

// @title AstroCat API
// @version 1.0
// @description AstroCat API sample for clients
// @BasePath /
func RunService(envSettings *schemas.EnvSettings, logger logging.Logger) {
	api, _ := NewApi(logger, envSettings)
	api.RunApi(envSettings)
}
