package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bekzat-kamen/RestApiFromZero/internal/database"
	"github.com/bekzat-kamen/RestApiFromZero/internal/handlers"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://taskuser:taskpass@localhost:5432/tasksdb?sslmode=disable"
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("Запуск сервера на порту %s", serverPort)

	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	log.Println("Успешное подключение к БД")

	taskStore := database.NewTaskStore(db)
	handler := handlers.NewHandlers(taskStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", methodHandler(handler.GetAllTasks, http.MethodGet))

	mux.HandleFunc("/tasks/create", methodHandler(handler.CreateTask, http.MethodPost))

	mux.HandleFunc("/tasks/", taskIDHandler(handler))
	loggedMux := loggingMiddleware(mux)

	serverAddr := ":" + serverPort
	log.Fatal(http.ListenAndServe(serverAddr, loggedMux))
}

func methodHandler(handlerFunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		handlerFunc(w, r)
	}
}

func taskIDHandler(handler *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTask(w, r)
		case http.MethodPut:
			handler.UpdateTask(w, r)
		case http.MethodDelete:
			handler.DeleteTask(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
