package temperatures

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func ListTemperatures(ctx *fiber.Ctx) error {
	cursor, err := temperaturesDb.Db.Collection("temperatures").Find(ctx.Context(), bson.D{{}})
	if err != nil {
		logrus.Error(err)
		return ctx.Status(500).JSON(fiber.Map{"message": "Error getting temperatures"})
	}
	var temperatures = make([]Temperature, 0)
	if err := cursor.All(ctx.Context(), &temperatures); err != nil {
		logrus.Error(err)
		return ctx.Status(500).JSON(fiber.Map{"message": "Error getting temperatures"})
	}
	return ctx.JSON(temperatures)
}
