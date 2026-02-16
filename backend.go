package main
import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

var (
	todos  = make(map[int]Todo)
	nextID = 1
	mu     sync.Mutex
)

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		return
	}

	mu.Lock()
	list := make([]Todo, 0)
	for _, todo := range todos {
		list = append(list, todo)
	}
	mu.Unlock()

	json.NewEncoder(w).Encode(list)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		return
	}
	var req struct {
		Title string `json:"title"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mu.Lock()
	id := nextID
	nextID++
	newTodo := Todo{
		ID:        id,
		Title:     req.Title,
		Done:      false,
		CreatedAt: time.Now(),
	}
	todos[id] = newTodo
	mu.Unlock()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		return
	}
	id, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/todos/"))
	var req struct {
		Done bool `json:"done"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	mu.Lock()
	if todo, exists := todos[id]; exists {
		todo.Done = req.Done
		todos[id] = todo
	}
	mu.Unlock()
	json.NewEncoder(w).Encode(todos[id])
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		return
	}
	id, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/todos/"))
	mu.Lock()
	delete(todos, id)
	mu.Unlock()
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/api/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getTodos(w, r)
		} else if r.Method == "POST" {
			createTodo(w, r)
		} else {
			enableCORS(w)
		}
	})
	http.HandleFunc("/api/todos/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			updateTodo(w, r)
		} else if r.Method == "DELETE" {
			deleteTodo(w, r)
		} else {
			enableCORS(w)
		}
	})
	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Frontend: http://localhost:8000")
	http.ListenAndServe(":8080", nil)
}
