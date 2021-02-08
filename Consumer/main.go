package main

import (
	svc "Consumer/service"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
)

var logger *log.Logger

func HttpKeepAlive(port string) {
	errChan := make(chan error)
	go func() {
		log.Println("HTTP KeepAlive :transport", "HTTP", "started on port", port)
		errChan <- http.ListenAndServe(port, nil)
	}()
	log.Fatal("exit", <-errChan)
}

func main() {
	port := GetValueFromEnvVariable("ENV_PORT", ":8080")
	producer := GetValueFromEnvVariable("PRODUCER_URL", "http://localhost:9090")
	consumer := svc.ConsumerService{
		StopChan: make(chan bool),
		Producer: producer,
		Redis: redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		}),
	}
	consumer.Initialize()
	HttpKeepAlive(port)
}
