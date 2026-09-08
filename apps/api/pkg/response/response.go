package response

import (
	"kun-galgame-sticker-api/pkg/errors"

	"github.com/gofiber/fiber/v3"
)

func OK(c fiber.Ctx, data any) error {
	return c.JSON(fiber.Map{
		"code":    errors.CodeOK,
		"message": "ok",
		"data":    data,
	})
}

// Error always emits data, as OK does. The two used to disagree -- success
// carried a data key and failure did not -- so a client could not decode both
// with one type.
func Error(c fiber.Ctx, err *errors.AppError) error {
	return c.Status(err.StatusCode).JSON(fiber.Map{
		"code":    err.Code,
		"message": err.Message,
		"data":    nil,
	})
}
