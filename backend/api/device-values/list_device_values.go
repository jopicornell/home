package device_values

import (
	"backend/log"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func ListDeviceValuesController(ctx *fiber.Ctx) error {
	logger := log.Logger.WithField("requestId", ctx.Locals("requestId"))
	cursor, err := deviceValuesCollection.Find(ctx.Context(), bson.D{{"_id", ctx.Params("device_id")}})
	if err != nil {
		logger.Error(err)
		return ctx.Status(500).JSON(fiber.Map{"message": "Error getting devices"})
	}
	var temperatures = make([]DeviceValue, 0)
	if err := cursor.All(ctx.Context(), &temperatures); err != nil {
		logger.Error(err)
		return ctx.Status(500).JSON(fiber.Map{"message": "Error getting devices"})
	}
	return ctx.JSON(temperatures)
}
