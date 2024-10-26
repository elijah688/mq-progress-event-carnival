package main

import (
	"encoding/json"
	"fmt"
	"log"
	"messages/internal/config"
	"messages/internal/middleware"
	"messages/internal/model"
	"messages/internal/queue"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
)

var messageMap sync.Map

func main() {
	cfg, err := config.NewQueueConfig()
	if err != nil {
		log.Fatal(err)
	}

	q, err := queue.NewQueue(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer q.Close()

	msgs, err := q.Consume()
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			var message model.Message
			if err := json.Unmarshal(msg.Body, &message); err != nil {
				log.Printf("Failed to parse message: %v", err)
				continue
			}

			messageMap.Store(message.ID, message)

			fmt.Printf("Received message ID: %s, Progress: %.2f\n", message.ID, message.PercentageComplete)
		}
	}()

	r := chi.NewRouter()

	r.Use(middleware.CORS)

	r.Get("/taskmon", taskmonHandler)
	r.Delete("/{id}", deleteTaskHandler)

	port := os.Getenv("CONSUMER_PORT")
	fmt.Printf("Server is running on port %s...\n", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func taskmonHandler(w http.ResponseWriter, r *http.Request) {
	tempMap := make(map[string]model.Message)

	messageMap.Range(func(key, value interface{}) bool {
		if msg, ok := value.(model.Message); ok {
			tempMap[key.(string)] = msg
		}
		return true // continue iteration
	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tempMap); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	if _, loaded := messageMap.LoadAndDelete(taskID); !loaded {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
