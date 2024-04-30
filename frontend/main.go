package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var connection *amqp.Connection

func main() {
	fmt.Println("Waiting for amqp")
	for {
		var err error
		connection, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err == nil {
			break
		} else {
			log.Println(err)
			time.Sleep(time.Second)
		}
	}

	mux := registerRoutes()

	fmt.Println("Frontend started on 8080")

	http.ListenAndServe(":8080", mux)
}
