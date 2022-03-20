package main

import "github.com/gofiber/fiber/v2"
import "backend/api/temperatures"

func main() {
	app := fiber.New()
	temperatures.TemperaturesApp(app)
	err := app.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
