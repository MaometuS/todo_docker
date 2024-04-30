package main

import (
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

var connection *amqp.Connection

func main() {
	mux := registerRoutes()

	var err error
	connection, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Frontend started on 8080")

	http.ListenAndServe(":8080", mux)
}
