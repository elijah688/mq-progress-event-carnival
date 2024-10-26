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
	"time"

	"github.com/google/uuid"
)

func main() {
	port := os.Getenv("PUBLISHER_PORT")

	cfg, err := config.NewQueueConfig()
	if err != nil {
		log.Fatal(err)
	}

	q, err := queue.NewQueue(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer q.Close()

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		taskID := uuid.New().String()
		go func() {

			message := model.Message{
				ID:                 taskID,
				Name:               "Sample Task",
				User:               "user@example.com",
				State:              "Running",
				StartTime:          time.Now().UTC().Format(time.RFC3339),
				FinishedTime:       time.Time{}.Format(time.RFC3339),
				Duration:           "0 secs",
				ErrorMessage:       "",
				PercentageComplete: 0.0,
			}

			for {

				if message.PercentageComplete > 1.0 {
					message.PercentageComplete = 1.0
					message.State = "Complete"
					message.FinishedTime = time.Now().UTC().Format(time.RFC3339)
					fmt.Printf("Completed task ID: %s\n", message.ID)
				}

				body, err := json.Marshal(message)
				if err != nil {
					log.Printf("Failed to marshal message: %v", err)
					return
				}

				if err := q.Publish(body); err != nil {
					log.Printf("Failed to publish message: %v", err)
					return
				}

				fmt.Printf("Published message ID: %s, Progress: %.2f, Publisher :%s\n", message.ID, message.PercentageComplete, port)

				if message.PercentageComplete == 1.0 {
					return
				}
				message.PercentageComplete += 0.01
				fmt.Println(message.PercentageComplete)

				time.Sleep(50 * time.Millisecond)
				fmt.Println(message.PercentageComplete)

			}

		}()

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("{\"taskID\":\"%s\"}", taskID))) // 202 Accepted
	})
	handler := middleware.CORS(http.DefaultServeMux)

	fmt.Printf("Server is running on port %s...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
