package main

import "net/http"

func registerRoutes() http.Handler {
	mux := &http.ServeMux{}

	mux.Handle("GET /", http.RedirectHandler("/todos", http.StatusMovedPermanently))

	mux.HandleFunc("GET /todos", getMainPage)
	mux.HandleFunc("POST /create_todo", addTodo)
	mux.HandleFunc("POST /edit_todo", editTodo)
	mux.HandleFunc("GET /delete_todo", removeTodo)

	return mux
}
