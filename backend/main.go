package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	fmt.Println("Waiting for amqp")
	for {
		var err error
		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err == nil {
			break
		} else {
			log.Println(err)
			time.Sleep(time.Second)
		}
	}

	fmt.Println("Waiting for db")
	for {
		var err error
		db, err = pgxpool.New(context.Background(), "postgres://postgres:postgres@postgres:5432/todo")
		if err == nil {
			break
		} else {
			log.Println(err)
			time.Sleep(time.Second)
		}
	}

	fmt.Println("Backend started")

	route()
}
