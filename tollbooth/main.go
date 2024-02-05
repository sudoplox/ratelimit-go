package main

import (
	"encoding/json"
	tollbooth "github.com/didip/tollbooth/v7"
	"log"
	"net/http"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func endpointHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK) // 200

	message := Message{
		Status: "Successful",
		Body:   "Hi you've reached the API",
	}
	err := json.NewEncoder(rw).Encode(&message)
	if err != nil {
		return
	}
}

func main() {
	message := Message{
		Status: "Request Failed",
		Body:   "The API is at capacity, please try again later.",
	}
	jsonMessage, _ := json.Marshal(message)
	tollBoothLimiter := tollbooth.NewLimiter(2, nil)
	tollBoothLimiter.SetMessageContentType("application/json")
	tollBoothLimiter.SetMessage(string(jsonMessage) + "\n")
	http.Handle("/ping", tollbooth.LimitFuncHandler(tollBoothLimiter, endpointHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("there was an error listening on port :8080", err)
	}
}
