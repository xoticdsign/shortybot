package middleware

import (
	"gopkg.in/telebot.v4"

	"github.com/xoticdsign/shortybot/internal/bot/models"
)

// Перехватывает все необходимые данные пользователя для последующего использования в хендлерах.
func GetUserDetails(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if c.Sender().Username == "" {
			return c.Send(models.FailedGlobalUsernameAbsent, models.ReplyReturnToMenuWithError)
		}
		c.Set("user", c.Sender())

		return next(c)
	}
}
