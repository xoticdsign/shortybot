package app

import (
	"os"
	"time"

	"gopkg.in/telebot.v4"

	"github.com/xoticdsign/shortybot/internal/bot/handlers"
	"github.com/xoticdsign/shortybot/internal/bot/helpers"
	"github.com/xoticdsign/shortybot/internal/bot/middleware"
	"github.com/xoticdsign/shortybot/internal/bot/models"
	"github.com/xoticdsign/shortybot/internal/db"
)

// Инициализирует бота, возвращает структуру *telebot.Bot или одну из возможных ошибок.
func InitApp() (*telebot.Bot, error) {
	db, err := db.InitDB()
	if err != nil {
		return nil, err
	}

	handlers := &handlers.Dependencies{
		DB:      db,
		Helpers: helpers.Helpers{},
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token: os.Getenv("BOT_TOKEN"),
		Poller: &telebot.LongPoller{
			Limit:   50,
			Timeout: time.Second * 15,
		},
		OnError: handlers.OnError,
	})
	if err != nil {
		return nil, err
	}

	bot.Use(middleware.GetUserDetails)

	unsupported := bot.Group()

	unsupported.Handle(telebot.OnText, handlers.Unsupported)
	unsupported.Handle(telebot.OnMedia, handlers.Unsupported)

	user := bot.Group()

	user.Handle(&models.BtnListShorties, handlers.ListShorties)
	user.Handle(&models.BtnShortyInfo, handlers.ShortyInfo)
	user.Handle(&models.BtnDeleteShorty, handlers.DeleteShorty)
	user.Handle(&models.BtnDeleteShortyPrompt, handlers.DeleteShortyPrompt)
	user.Handle(&models.BtnDeleteSelectedShorty, handlers.DeleteSelectedShorty)
	user.Handle(&models.BtnReturnToMenu, handlers.Menu)

	return bot, nil
}