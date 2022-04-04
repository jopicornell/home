package devices

import (
	"backend/log"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func ListDevicesController(ctx *fiber.Ctx) error {
	logger := log.Logger.WithField("requestId", ctx.Locals("requestId"))
	cursor, err := devicesCollection.Find(ctx.Context(), bson.D{})
	if err != nil {
		logger.Error(err)
		return ctx.Status(500).JSON(fiber.Map{"message": "Error getting devices"})
	}
	var devices = make([]Device, 0)
	if err := cursor.All(ctx.Context(), &devices); err != nil {
		logger.Error(err)
		return ctx.Status(500).JSON(fiber.Map{"message": "Error getting devices"})
	}
	return ctx.JSON(devices)
}
