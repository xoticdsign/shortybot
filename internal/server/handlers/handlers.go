package handlers

import (
	"errors"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"

	"github.com/xoticdsign/shortybot/internal/db"
	"github.com/xoticdsign/shortybot/internal/server/models"
)

// Структура, хранящая все необходимые хендлерам сервера зависимости.
type Dependencies struct {
	DB db.Querier
}

// Отлавливает ошибки и обрабатывает.
func (d *Dependencies) OnError(c *fiber.Ctx, err error) error {
	var e *fiber.Error

	if errors.As(err, &e) {
		return c.JSON(e)
	}
	return c.JSON(models.Error{
		Code:    0,
		Message: err.Error(),
	})
}

// Достает оригинальную ссылку из БД и редиректит запросы.
func (d *Dependencies) Redirect(c *fiber.Ctx) error {
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
	return c.Redirect(shorty.OriginalURL, fiber.StatusSeeOther)
}
