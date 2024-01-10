package rest

import (
	"fmt"

	"github.com/RuanScherer/journey-track-api/application/rest"
	"github.com/RuanScherer/journey-track-api/application/rest/middlewares"
	"github.com/RuanScherer/journey-track-api/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func StartAPI() {
	appConfig := config.GetAppConfig()
	app := fiber.New(fiber.Config{
		AppName:      "Journey Track API",
		ErrorHandler: middlewares.HandleError,
	})
	app.Use(logger.New())
	rest.RegisterRoutes(app)
	app.Listen(fmt.Sprintf(":%v", appConfig.RestApiPort))
}
