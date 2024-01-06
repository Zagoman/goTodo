package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	listenAddr string
	store      Storage
}

func NewApiServer(listenAddr string, store Storage) *ApiServer {
	return &ApiServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/todo", makeHTTPHandlerFunc(s.handleTodo))
	router.HandleFunc("/todo/{id}", makeHTTPHandlerFunc(s.handleUniqueTodo))
	log.Println("JSON API listening on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *ApiServer) handleCreateTodo(w http.ResponseWriter, r *http.Request) error {
	createTodoRequest := &CreateTodoRequest{}
	if err := json.NewDecoder(r.Body).Decode(createTodoRequest); err != nil {
		return err
	}

	todo := &Todo{
		Task: createTodoRequest.Task,
	}
	task, err := s.store.CreateTodo(todo)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, task)
}
func (s *ApiServer) handleDeleteTodo(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, &APIError{Error: "The ID must be an integer value"})
	}
	todo, err := s.store.DeleteTodo(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, todo)
}
func (s *ApiServer) handleUpdateTodo(w http.ResponseWriter, r *http.Request) error {
	updateTodoRequest := &UpdateTodoRequest{}
	if err := json.NewDecoder(r.Body).Decode(updateTodoRequest); err != nil {
		return err
	}

	task := &Todo{
		Task: updateTodoRequest.Task,
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, &APIError{Error: "The ID must be an integer value"})
	}
	todo, err := s.store.UpdateTodo(id, task)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, todo)
}
func (s *ApiServer) handleGetTodoById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, &APIError{Error: "The ID must be an integer value"})
	}
	todo, err := s.store.GetTodoById(id)
	if err != nil {
		return err
	}
	fmt.Println(todo)
	return WriteJSON(w, http.StatusOK, todo)
}
func (s *ApiServer) handleGetTodo(w http.ResponseWriter, r *http.Request) error {
	todos, err := s.store.GetTodos()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, todos)
}
func (s *ApiServer) handleTodo(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetTodo(w, r)
	} else if r.Method == "POST" {
		return s.handleCreateTodo(w, r)
	}
	return WriteJSON(w, http.StatusMethodNotAllowed, &APIError{Error: "Method not allowed"})
}

func (s *ApiServer) handleUniqueTodo(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "DELETE" {
		return s.handleDeleteTodo(w, r)
	} else if r.Method == "PATCH" {
		return s.handleUpdateTodo(w, r)
	} else if r.Method == "GET" {
		return s.handleGetTodoById(w, r)
	}
	return WriteJSON(w, http.StatusMethodNotAllowed, &APIError{Error: "Method not allowed"})
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle error here
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}
