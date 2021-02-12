package client

import (
	v1 "Producer/contracts"
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var taskList []v1.Task

// Fetch the environment variable if it exists or else return the defaultValue set in the code.
func GetValueFromEnvVariable(variableName, defaultValue string) string {
	environmentValue := os.Getenv(variableName)
	if environmentValue == "" {
		return defaultValue
	}
	return environmentValue
}

/* Can be implemented as per client*/
func ProducerClient() {
	client := http.Client{}
	taskCount := GetValueFromEnvVariable("TASK_COUNT", "10")
	count, err := strconv.Atoi(taskCount)
	if err != nil {
		log.Print(err)
	}
	for i := 0; i < count; i++ {
		reqBodyBytes := new(bytes.Buffer)
		task := BuildRandomRequest(count)
		json.NewEncoder(reqBodyBytes).Encode(task)
		consumer := GetValueFromEnvVariable("PRODUCER_URL", "http://localhost:9090")
		consumerEndPoint := "/tasks/produce"
		request, err := http.NewRequest("POST", consumer+consumerEndPoint, bytes.NewBuffer(reqBodyBytes.Bytes()))
		if err != nil {
			log.Fatal("Unable to POST task to Producer")
		}

		_, err = client.Do(request)
		if err != nil {
			log.Fatal("Error Producing task")
		}
	}
}

func randomDate(count int) time.Time {
	start := time.Date(2021, 2, 11, 12, 0, 0, 0, time.UTC)
	randomDate := start.Add(time.Hour * time.Duration(count))
	return randomDate
}

func BuildRandomRequest(count int) v1.Task {
	taskRequest := v1.Task{}
	random := rand.Intn(count)
	taskRequest.TaskName = "task_" + string(random)
	taskRequest.TaskStatus = "pending"
	taskRequest.Periodicity = 5
	taskRequest.ScheduledTime = randomDate(count).String()
	taskRequest.LastUpdateTime = time.Now().UTC().String()
	taskRequest.TaskType = "task_test"
	return taskRequest
}
