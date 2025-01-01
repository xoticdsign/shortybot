package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/xoticdsign/shortybot/internal/bot/app"
)

func main() {
	if os.Getenv("ENV") != "production" {
		godotenv.Load()
	}

	bot, err := app.InitApp()
	if err != nil {
		log.Fatalf("ERR %v", err)
	}

	bot.Start()
}
