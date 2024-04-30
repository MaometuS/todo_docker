package main

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rpcRequest struct {
	Method string
	Data   any
}

type todo struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func makeCall(req rpcRequest) ([]byte, error) {
	chn, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	q, err := chn.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := chn.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	corrID := uuid.New().String()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	err = chn.PublishWithContext(
		ctx,
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrID,
			ReplyTo:       q.Name,
			Body:          body,
		},
	)
	if err != nil {
		return nil, err
	}

	for d := range msgs {
		if corrID == d.CorrelationId {
			return d.Body, nil
		}
	}

	return nil, errors.New("could not receive body")
}

func getTodos() ([]todo, error) {
	body, err := makeCall(rpcRequest{Method: "get_todos"})
	if err != nil {
		return nil, err
	}

	todos := make([]todo, 0)
	err = json.Unmarshal(body, &todos)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func createTodo(name string) error {
	_, err := makeCall(rpcRequest{Method: "create_todo", Data: todo{Name: name}})
	if err != nil {
		return err
	}

	return nil
}

func updateTodo(t *todo) error {
	_, err := makeCall(rpcRequest{Method: "update_todo", Data: t})
	if err != nil {
		return err
	}

	return nil
}

func deleteTodo(id int64) error {
	_, err := makeCall(rpcRequest{Method: "delete_todo", Data: todo{ID: id}})
	if err != nil {
		return err
	}

	return nil
}
