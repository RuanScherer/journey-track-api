package rest

import (
	"fmt"

	"github.com/RuanScherer/journey-track-api/adapters/rest/middlewares"
	"github.com/RuanScherer/journey-track-api/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func StartServer() {
	appConfig := config.GetAppConfig()
	app := fiber.New(fiber.Config{
		AppName:      "Journey Track API",
		ErrorHandler: middlewares.HandleError,
	})
	app.Use(logger.New())
	RegisterRoutes(app)
	app.Listen(fmt.Sprintf(":%v", appConfig.RestApiPort))
}
