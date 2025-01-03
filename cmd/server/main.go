package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/xoticdsign/shortybot/internal/server/app"
)

func main() {
	env := os.Getenv("ENV")

	defer os.Unsetenv("ENV")

	if env != "production" {
		godotenv.Load()
	}

	app, err := app.InitApp()
	if err != nil {
		log.Fatalf(" ERR %v", err)
	}

	err = app.Listen(os.Getenv("SERVER_ADR"))
	if err != nil {
		log.Fatalf(" ERR %v", err)
	}
}
