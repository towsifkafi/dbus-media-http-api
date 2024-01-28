package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	setupLogging()

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default port [:10004] and no auth")
	}

	handleRequests()
}
