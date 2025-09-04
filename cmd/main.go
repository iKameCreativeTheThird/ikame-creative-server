package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env file")
	}
	port := os.Getenv("PORT")
	fmt.Printf("Hello, Performance Dashboard! Running on port %s\n", port)
}
