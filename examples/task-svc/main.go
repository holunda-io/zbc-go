package main

import (
	"encoding/json"
	"fmt"
	"github.com/zeebe-io/zbc-go/zbc"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	ProcessedEventsCount uint64
	ErrorCount           uint64
)

const BrokerAddr = "0.0.0.0:51015"

type StopCh chan bool

var workers map[string]StopCh

func processTask(lo string, msg *zbc.TaskEvent) {
	log.Printf("[%s] Working on task.\n", lo)
}

func openSubscription(client *zbc.Client, stopCh chan bool, topic string, lo string, tt string) {
	subscriptionCh, subInfo, err := client.TaskConsumer(topic, lo, tt)
	if err != nil {
		atomic.AddUint64(&ErrorCount, 1)
	}

	credits := subInfo.Credits

	log.Println("Subscription opened with", credits, "Credits")
	log.Println("Waiting for events ....")
	for {
		select {
		case message := <-subscriptionCh:
			credits--;

			processTask(lo, message)
			response, err := client.CompleteTask(message)

			if err != nil {
				log.Println("Completing a task went wrong.")
				log.Println(err)
			}

			if response.State == zbc.TaskCompleted {
				atomic.AddUint64(&ProcessedEventsCount, 1)
				log.Println("Task completed successfully.")
			} else {
				log.Println("Task not completed.")
			}

			if credits < 1 {
				response, err := client.IncreaseTaskSubscriptionCredits(subInfo)

				if err != nil {
					log.Println("Increasing task credits went wrong.")
					log.Println(err)

				} else {
					credits = response.Credits
					log.Println("Increased task credits to", credits)
				}
			}

			break

		case stop := <-stopCh:
			if stop {
				log.Print("Stopping worker.")
				_, err := client.CloseTaskSubscription(subInfo)
				if err != nil {
					log.Println("Close task subscription request failed")
					log.Println(err)
				}
				log.Println("Gracefully shutting down the client.")
				client.Close()
				return
			}
			break
		}
	}
}

func startWorkerView(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	// TODO: fix multiple clients ...
	zbClient, err := zbc.NewClient(BrokerAddr)
	if err != nil {
		resp["status"] = http.StatusInternalServerError
	}
	resp["status"] = http.StatusOK
	lockOwner := fmt.Sprintf("zbc-%s", time.Now().Format("20060102150405"))

	log.Printf("Starting worker with ID: %s\n", lockOwner)
	workers[lockOwner] = make(chan bool)
	go openSubscription(zbClient, workers[lockOwner], "default-topic", lockOwner, "foo")

	resp["workerID"] = lockOwner
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(jsonResp)
}

func statsView(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	resp["Running"] = true
	if len(workers) == 0 {
		resp["Running"] = false
	} else {
		resp["WorkersRunning"] = len(workers)
		resp["ProcessedEventsCount"] = atomic.LoadUint64(&ProcessedEventsCount)
		resp["ErrorCount"] = atomic.LoadUint64(&ErrorCount)

		var workersIDs []string
		for key, _ := range workers {
			workersIDs = append(workersIDs, key)
		}
		resp["WorkersIDs"] = workersIDs
	}
	jsonResp, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(jsonResp)
}

func stopWorkersView(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})
	for _, worker := range workers {
		worker <- true
	}
	workers = make(map[string]StopCh)
	resp["Status"] = http.StatusOK

	jsonResp, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(jsonResp)
}

func main() {
	workers = make(map[string]StopCh)
	log.Println("Super microservice started.")
	log.Println("Waiting for workers to start.")

	http.HandleFunc("/start", startWorkerView)
	http.HandleFunc("/stop", stopWorkersView)
	http.HandleFunc("/stats", statsView)

	http.ListenAndServe(":3000", nil)
}
