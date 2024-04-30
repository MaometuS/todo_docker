package main

import (
	"encoding/json"
	"time"
)

type todo struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func getTodos(req []byte) ([]byte, error) {
	todos := []todo{
		{
			ID:   1,
			Name: "name1",
		},
		{
			ID:   2,
			Name: "name2",
		},
	}

	return json.Marshal(todos)
}

func createTodo(req []byte) ([]byte, error) {
	panic("impl")
}

func updateTodo(req []byte) ([]byte, error) {
	panic("impl")
}

func deleteTodo(req []byte) ([]byte, error) {
	panic("impl")
}
