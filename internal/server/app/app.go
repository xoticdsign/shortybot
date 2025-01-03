package app

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/xoticdsign/shortybot/internal/db"
	"github.com/xoticdsign/shortybot/internal/logger"
	"github.com/xoticdsign/shortybot/internal/server/handlers"
)

// Инициализирует сервер, возвращает структуру *fiber.App или одну из возможных ошибок.
func InitApp() (*fiber.App, error) {
	dsn := os.Getenv("DB_DSN")

	defer os.Unsetenv("DB_DSN")

	db, err := db.InitDB(dsn)
	if err != nil {
		return nil, err
	}

	logger := logger.InitLogger()

	handlers := &handlers.Dependencies{
		DB:     db,
		Logger: logger,
	}

	app := fiber.New(fiber.Config{
		ServerHeader:  "shortyserver",
		StrictRouting: true,
		CaseSensitive: true,
		ReadTimeout:   time.Second * 20,
		WriteTimeout:  time.Second * 20,
		ErrorHandler:  handlers.OnError,
		AppName:       "shortyserver",
	})

	app.Get("/", handlers.Root)
	app.Get("/:shortyURL", handlers.Redirect)

	return app, nil
}
