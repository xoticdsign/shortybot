package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/xoticdsign/shortybot/internal/bot/app"
)

func main() {
	env := os.Getenv("ENV")

	defer os.Unsetenv("ENV")

	if env != "production" {
		godotenv.Load()
	}

	bot, err := app.InitApp()
	if err != nil {
		log.Fatalf("ERR %v", err)
	}

	bot.Start()
}
