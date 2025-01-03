package middleware

import (
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/telebot.v4"
)

// Перехватывает все необходимые данные пользователя для последующего использования в хендлерах, а также проверяет, является ли пользователем админом.
func GetSenderDetails(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		c.Set("user", c.Sender())

		return next(c)
	}
}

// Проверяет, является ли пользователь админом.
func AdminValidation(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		admins := os.Getenv("BOT_ADMINS")
		if !strings.Contains(admins, strconv.Itoa(int(c.Sender().ID))) {
			return next(c)
		}
		c.Set("admin", "")

		return next(c)
	}
}

// Счетчик скорости.
func SpeedCounter(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		c.Set("start", time.Now())

		return next(c)
	}
}
