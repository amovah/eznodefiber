package eznodefiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type customErrorResponse struct {
	Message string `json:"message"`
}

func errorHandler(logger *logrus.Logger) func(ctx *fiber.Ctx, err error) error {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		message := "Internal Server Error"

		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			message = e.Message
		}

		if code >= 500 {
			logger.Errorln(err)
		} else {
			logger.Warnln(err)
		}

		err = ctx.Status(code).JSON(customErrorResponse{
			Message: message,
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		return nil
	}
}
