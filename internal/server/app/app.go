package app

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/xoticdsign/shortybot/internal/db"
	"github.com/xoticdsign/shortybot/internal/logger"
	"github.com/xoticdsign/shortybot/internal/server/handlers"
)

// Инициализирует сервер, возвращает структуру *fiber.App или одну из возможных ошибок.
func InitApp() (*fiber.App, error) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDb := os.Getenv("DB_DB")
	dbPort := os.Getenv("DB_PORT")

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
