package app

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/telebot.v4"

	"github.com/xoticdsign/shortybot/internal/bot/handlers"
	"github.com/xoticdsign/shortybot/internal/bot/helpers"
	"github.com/xoticdsign/shortybot/internal/bot/middleware"
	"github.com/xoticdsign/shortybot/internal/bot/models"
	"github.com/xoticdsign/shortybot/internal/db"
	"github.com/xoticdsign/shortybot/internal/logger"
)

// Инициализирует бота, возвращает структуру *telebot.Bot или одну из возможных ошибок.
func InitApp() (*telebot.Bot, error) {
	token := os.Getenv("BOT_TOKEN")
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDb := os.Getenv("DB_DB")
	dbPort := os.Getenv("DB_PORT")

	defer os.Unsetenv("BOT_TOKEN")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_USER")
	defer os.Unsetenv("DB_PASSWORD")
	defer os.Unsetenv("DB_DB")
	defer os.Unsetenv("DB_PORT")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v", dbHost, dbUser, dbPassword, dbDb, dbPort)

	db, err := db.InitDB(dsn)
	if err != nil {
		return nil, err
	}

	logger := logger.InitLogger()

	handlers := &handlers.Dependencies{
		DB:      db,
		Logger:  logger,
		Helpers: helpers.Helpers{},
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token: token,
		Poller: &telebot.LongPoller{
			Limit:   50,
			Timeout: time.Second * 15,
		},
		OnError: handlers.OnError,
	})
	if err != nil {
		return nil, err
	}

	bot.Use(middleware.GetSenderDetails)
	bot.Use(middleware.AdminValidation)
	bot.Use(middleware.SpeedCounter)

	admin := bot.Group()

	admin.Handle(&models.BtnReturnToAdminPanel, handlers.AdminPanel)
	admin.Handle(&models.BtnAdminUsersAndShorties, handlers.AdminUsersAndShorties)

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
