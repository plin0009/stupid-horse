package main

import (
	"os"

	"github.com/joho/godotenv"
)

func main() {
}

func startBot() {
	err := godotenv.Load()
	if err != nil {
		panic("could not load .env file")
	}
	b := Bot{token: os.Getenv("lichess_key")}
	b.Listen()
}
