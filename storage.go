package main

type Storage interface {
	CreateTodo(string) error
	DeleteTodo(string) error
	GetTodoById(int) (*Todo, error)
	GetTodos() ([]*Todo, error)
}
