package main

import (
	svc "Producer/service"
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
	port := GetValueFromEnvVariable("ENV_PORT", ":9090")
	consumer := GetValueFromEnvVariable("CONSUMER_URL", "http://localhost:8080")
	producer := svc.ProducerService{
		StopChan: make(chan bool),
		Consumer: consumer,
	}
	producer.Initialize()
	HttpKeepAlive(port)
}
