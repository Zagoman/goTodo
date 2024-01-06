package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateTodo(*Todo) (*Todo, error)
	DeleteTodo(int) (*Todo, error)
	UpdateTodo(int, *Todo) (*Todo, error)
	GetTodoById(int) (*Todo, error)
	GetTodos() ([]*Todo, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gotodo sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateTodoTable()
}

func (s *PostgresStore) CreateTodoTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS todo (
		id SERIAL PRIMARY KEY,
		task VARCHAR(255),
		created_at TIMESTAMP
	)
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateTodo(task *Todo) (*Todo, error) {

	todo := &Todo{}

	stmt, err := s.db.Prepare(`INSERT INTO todo 
	(task, created_at)
	VALUES	
	($1, $2)
	RETURNING * 
	`)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRow(task.Task, task.CreatedAt).Scan(&todo.ID, &todo.Task, &todo.CreatedAt)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *PostgresStore) DeleteTodo(id int) (*Todo, error) {
	todo := &Todo{}

	stmt, err := s.db.Prepare(`
		DELETE FROM todo WHERE id=$1 RETURNING *	
	`)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRow(id).Scan(&todo.ID, &todo.Task, &todo.CreatedAt)
	if err != nil {
		return nil, err
	}
	return todo, nil
}
func (s *PostgresStore) GetTodoById(id int) (*Todo, error) {
	todo := &Todo{}
	row := s.db.QueryRow(`
	SELECT * FROM todo 
	WHERE id=$1	`, id)
	if err := row.Scan(&todo.ID, &todo.Task, &todo.CreatedAt); err != nil {
		return nil, err
	}

	return todo, nil
}
func (s *PostgresStore) UpdateTodo(id int, task *Todo) (*Todo, error) {
	todo := &Todo{}

	stmt, err := s.db.Prepare(`
		UPDATE todo 
		SET task=$1
		WHERE id=$2
		RETURNING *
	`)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRow(task.Task, id).Scan(&todo.ID, &todo.Task, &todo.CreatedAt)
	if err != nil {
		return nil, err
	}
	return todo, nil
}
func (s *PostgresStore) GetTodos() ([]*Todo, error) {
	rows, err := s.db.Query("SELECT * FROM todo")
	if err != nil {
		return nil, err
	}
	todos := []*Todo{}
	for rows.Next() {
		todo := &Todo{}
		err := rows.Scan(
			&todo.ID,
			&todo.Task,
			&todo.CreatedAt)
		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}
	return todos, nil
}
