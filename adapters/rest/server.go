package rest

import (
	"fmt"

	"github.com/RuanScherer/journey-track-api/adapters/rest/middleware"
	"github.com/RuanScherer/journey-track-api/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func StartServer() {
	appConfig := config.GetAppConfig()
	app := fiber.New(fiber.Config{
		AppName:      "Journey Track API",
		ErrorHandler: middleware.HandleError,
	})

	app.Use(logger.New())
	if appConfig.Environment == "development" {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     "http://localhost:3000",
			AllowCredentials: true,
		}))
	}

	RegisterRoutes(app)
	app.Listen(fmt.Sprintf(":%v", appConfig.RestApiPort))
}
