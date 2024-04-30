package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection

type rpc_request struct {
	Method string
	Data   json.RawMessage
}

func handleDelivery(d amqp.Delivery, ctx context.Context) {
	chn, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer chn.Close()

	req := new(rpc_request)
	err = json.Unmarshal(d.Body, req)
	if err != nil {
		err = chn.PublishWithContext(ctx,
			"",        // exchange
			d.ReplyTo, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId,
				Body:          []byte(err.Error()),
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		d.Ack(false)
		return
	}

	var res []byte
	switch req.Method {
	case "get_todos":
		res, err = getTodos(req.Data)
	case "create_todo":
		res, err = createTodo(req.Data)
	case "update_todo":
		res, err = updateTodo(req.Data)
	case "delete_todo":
		res, err = deleteTodo(req.Data)
	}

	if err != nil {
		err = chn.PublishWithContext(ctx,
			"",        // exchange
			d.ReplyTo, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId,
				Body:          []byte(err.Error()),
			},
		)
		if err != nil {
			log.Fatal(err)
		}

		d.Ack(false)
		return
	}

	err = chn.PublishWithContext(ctx,
		"",        // exchange
		d.ReplyTo, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: d.CorrelationId,
			Body:          res,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	d.Ack(false)
}

func route() {
	chn, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer chn.Close()

	q, err := chn.QueueDeclare(
		"rpc_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := chn.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for d := range msgs {
		go handleDelivery(d, ctx)
	}
}
