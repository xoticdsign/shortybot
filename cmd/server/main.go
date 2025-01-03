package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/xoticdsign/shortybot/internal/server/app"
)

func main() {
	env := os.Getenv("ENV")
	adr := os.Getenv("SERVER_ADR")

	defer os.Unsetenv("ENV")
	defer os.Unsetenv("SERVER_ADR")

	if env != "production" {
		godotenv.Load()
	}

	app, err := app.InitApp()
	if err != nil {
		log.Fatalf(" ERR %v", err)
	}

	err = app.Listen(adr)
	if err != nil {
		log.Fatalf(" ERR %v", err)
	}
}
