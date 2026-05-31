package main

import (
	"api-handler/db"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	conn, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Connected to PGDB")
	fmt.Println("Starting server at port :8080")

	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("Cannot start server at port :8080")
	}
}
