package main

import (
	"backend/api/device-values"
	"backend/api/devices"
	"backend/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

func main() {
	log.CreateLoggers()
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("request_id", uuid.New().String())
		return c.Next()
	})
	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path} ${body}\n",
	}))
	api := app.Group("/api")
	api.Mount("", device_values.App())
	api.Mount("", devices.App())

	err := app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
