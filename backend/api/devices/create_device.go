package devices

import (
	"backend/lib"
	"backend/log"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateDeviceController(ctx *fiber.Ctx) error {
	var device Device
	logger := log.Logger.WithField("requestId", ctx.Locals("request_id"))
	err := ctx.BodyParser(&device)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	err = device.Validate()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(lib.ValidationErrorsJson(err))
	}
	result, err := device.GetCollection().InsertOne(ctx.Context(), device)
	if err != nil {
		logger.Error(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error inserting document",
		})
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		_ = ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error getting inserted id",
		})
		logger.Error(errors.New("Error getting inserted id"))
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error inserting document",
		})
	}
	device.SetId(oid)
	logrus.Info("Inserted document with id: ", oid.Hex())
	ctx.Append("Location", fmt.Sprintf("%s/%s", ctx.Path(), oid.Hex()))
	return ctx.Status(fiber.StatusCreated).JSON(device)
}
