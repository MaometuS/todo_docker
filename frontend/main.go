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
	connection, err = amqp.Dial("")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Started on 8080")

	http.ListenAndServe(":8080", mux)
}
