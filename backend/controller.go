package main

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var db *pgxpool.Pool

type todo struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func getTodos(req []byte) ([]byte, error) {
	rows, err := db.Query(context.Background(), "select * from todos")
	if err != nil {
		return nil, err
	}

	todos, err := pgx.CollectRows(rows, pgx.RowToStructByName[todo])
	if err != nil {
		return nil, err
	}

	return json.Marshal(todos)
}

func createTodo(req []byte) ([]byte, error) {
	td := new(todo)
	err := json.Unmarshal(req, td)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(context.Background(), "insert into todos(name, created_at, updated_at) values ($1, $2, $3)", td.Name, time.Now(), time.Now())
	if err != nil {
		return nil, err
	}

	return []byte("success"), nil
}

func updateTodo(req []byte) ([]byte, error) {
	td := new(todo)
	err := json.Unmarshal(req, td)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(context.Background(), "update todos set name=$1, updated_at=$2 where id=$3", td.Name, time.Now(), td.ID)
	if err != nil {
		return nil, err
	}

	return []byte("success"), nil
}

func deleteTodo(req []byte) ([]byte, error) {
	td := new(todo)
	err := json.Unmarshal(req, td)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(context.Background(), "delete from todos where id=$1", td.ID)
	if err != nil {
		return nil, err
	}

	return []byte("success"), nil
}
