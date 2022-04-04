package device_values

import (
	"backend/lib"
	"github.com/gofiber/fiber/v2"
)

func CreateDeviceValueController(ctx *fiber.Ctx) error {
	var deviceValue DeviceValue
	err := lib.InsertFromBody(ctx, &deviceValue)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(deviceValue)
}
