package service

import (
	v1 "Producer/contracts"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type ProducerService struct {
	Collector Collector
	Router    *mux.Router
	StopChan  chan bool
	Consumer  string
}

func (prd *ProducerService) Initialize() {
	prd.Collector = NewCollector()
	http.HandleFunc("/producer/start", prd.InitializeDispatcher)
	http.HandleFunc("/producer/stop", prd.StopProducer)
	http.HandleFunc("/tasks/produce", prd.Collector.RequestCollector)
}

// Function to capture the task-requests in buffered channels and produce to Consumer
func (prd *ProducerService) InitializeDispatcher(_ http.ResponseWriter, _ *http.Request) {
	log.Print("Producer Started")
	go prd.StartProducer()
}

func (prd *ProducerService) StartProducer() {
	for {
		select {
		case task := <-TaskChan:
			log.Print("task", task.TaskName, "ready to be produced")
			// if successful then do nothing print success
			response, err := prd.ConsumerClient(task)
			if err != nil || response == "failure" {
				TaskChan <- task
				log.Print("task", task.TaskName, "failed to be consumed", task.TaskName, "enqueued")
			}
			log.Print("task", task.TaskName, "consumed successfully")
		case <-prd.StopChan:
			log.Print("Producer Stopped")
			break
		}
	}
}

func (prd *ProducerService) StopProducer(_ http.ResponseWriter, _ *http.Request) {
	go func() {
		prd.StopChan <- true
	}()
}

func (prd *ProducerService) ConsumerClient(task v1.Task) (string, error) {
	client := http.Client{}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(task)
	consumer := prd.Consumer
	consumerEndPoint := "/tasks/consume"
	request, err := http.NewRequest("POST", consumer+consumerEndPoint, bytes.NewBuffer(reqBodyBytes.Bytes()))
	if err != nil {
		log.Fatal("Unable to POST task to consumer")
		return "", err
	}

	response, err := client.Do(request)
	if err != nil {
		log.Fatal("Unexpected response from consumer")
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Unexpected body from consumer")
		return "", err
	}
	return string(body), nil
}
