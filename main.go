package main

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/raysandeep/Estimator-App/api/router"
	"github.com/raysandeep/Estimator-App/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func healthCheck(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func main() {
	utils.ImportEnv()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())
	app.Get("/", healthCheck)
	router.MountRoutes(app)

	app.Listen(":" + viper.GetString("PORT"))
}
