package temperatures

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTemperature(ctx *fiber.Ctx) error {
	// decode fiber request body to temperature
	var temperature Temperature
	err := ctx.BodyParser(&temperature)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err != nil {
		logrus.Error(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to generate uuid for temperature",
		})
	}
	result, err := mongoCollection.InsertOne(ctx.Context(), temperature)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error inserting temperature",
		})
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error converting temperature id to object id",
		})
	}
	logrus.Info("Inserted temperature with id: ", oid.Hex())
	ctx.Append("Location", "/temperatures/"+oid.Hex())
	return ctx.Status(fiber.StatusCreated).JSON(temperature)
}
