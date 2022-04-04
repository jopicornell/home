package lib

import "github.com/gofiber/fiber/v2"

func MongoBodyParser(ctx *fiber.Ctx, model MongoModel) error {
	// Parse JSON body
	err := ctx.BodyParser(&model)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return nil
}
