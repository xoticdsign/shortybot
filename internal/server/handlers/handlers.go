package handlers

import (
	"errors"
	"os"
	"time"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"

	"github.com/xoticdsign/shortybot/internal/db"
	"github.com/xoticdsign/shortybot/internal/logger"
	"github.com/xoticdsign/shortybot/internal/server/models"
)

// Структура, хранящая все необходимые хендлерам сервера зависимости.
type Dependencies struct {
	DB     db.Querier
	Logger logger.Loggier
}

// Отлавливает ошибки и обрабатывает.
func (d *Dependencies) OnError(c *fiber.Ctx, err error) error {
	var e *fiber.Error

	if errors.As(err, &e) {
		d.Logger.ErrorServer(
			e.Error(),
			e.Code,
		)

		return c.JSON(e)
	}
	d.Logger.ErrorServer(
		err.Error(),
		0,
	)

	return c.JSON(models.Error{
		Code:    0,
		Message: err.Error(),
	})
}

// Root хендлер.
func (d *Dependencies) Root(c *fiber.Ctx) error {
	return c.JSON(models.Error{
		Code:    200,
		Message: "OK",
	})
}

// Достает оригинальную ссылку из БД и редиректит запросы.
func (d *Dependencies) Redirect(c *fiber.Ctx) error {
	start := time.Now()

	shortyURL := c.Params("shortyURL")

	shorty, err := d.DB.ShortyInfo(shortyURL)
	if err != nil {
		switch {
		case err == gorm.ErrRecordNotFound:
			return fiber.ErrNotFound

		default:
			return err
		}
	}
	d.Logger.InfoServer(
		logger.InfoRequestFulfilled,
		os.Getenv("SERVER_PUBLIC_ADR")+shorty.ShortyURL,
		shorty.OriginalURL,
		start,
	)

	return c.Redirect(shorty.OriginalURL, fiber.StatusSeeOther)
}
