package main

import (
	"math/rand"
	"time"
)

type CreateTodoRequest struct {
	Task string `json:"task"`
}

type UpdateTodoRequest = CreateTodoRequest
type Todo struct {
	ID        int       `json:"id"`
	Task      string    `json:"task"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewTodo(task string) *Todo {
	return &Todo{
		ID:        rand.Intn(1000),
		Task:      task,
		CreatedAt: time.Now().UTC(),
	}
}
