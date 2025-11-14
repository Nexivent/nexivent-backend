package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "onichankimochi.com/astro_cat_backend/src/server/api/docs" // Import generated swagger docs
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// HealthCheck 			godoc
// @Summary 			Health Check
// @Description 		Verify connection in swagger
// @Tags 				Health Check
// @Accept 				json
// @Produce 			json
// @Success 			200 {object} 	string "ok"
// @Router 				/health-check/ [get]
func (a *Api) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "Works well !!")
}

func (a *Api) RunApi(envSettings *schemas.EnvSettings) {
	// CORS
	corsConfig := middleware.CORSConfig{
		AllowOrigins:     []string{"*"}, // TODO: allow only FrontOffice and UserFront origins
		AllowCredentials: true,
	}
	a.Echo.Use(middleware.CORSWithConfig(corsConfig))

	if envSettings.EnableSwagger {
		a.Echo.GET("/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.InstanceName("server")))
	}

	// Start the server
	a.Logger.Infoln(fmt.Sprintf("AstroCat server running on port %s", a.EnvSettings.MainPort))
	a.Logger.Fatal(a.Echo.Start(fmt.Sprintf(":%s", a.EnvSettings.MainPort)))
}
