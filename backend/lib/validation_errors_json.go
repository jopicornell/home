package lib

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidationErrorsJson(errors error) fiber.Map {
	errs, ok := errors.(validator.ValidationErrors)
	if !ok {
		return nil
	}
	fiberMap := fiber.Map{}
	for _, err := range errs {
		fiberMap[err.Field()] = fiber.Map{
			"namespace": err.Namespace(),
			"field":     err.Field(),
			"tag":       err.Tag(),
			"actualTag": err.ActualTag(),
			"param":     err.Param(),
			"error":     err.Error(),
		}
	}
	return fiberMap
}
