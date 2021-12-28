package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raysandeep/Estimator-App/api/controllers"
)

// MountRoutes mounts all routes declared here
func MountRoutes(app *fiber.App) {
	app.Get("/api/token/:channel", controllers.CreateRTCToken)
	app.Post("/api/start/call", controllers.StartCall)
	app.Post("/api/stop/call", controllers.StopCall)
	app.Get("/download", controllers.DownloadVideo)
	app.Get("/play", controllers.PlayVideo)
}
