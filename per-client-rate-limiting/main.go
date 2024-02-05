package main

import (
	"encoding/json"
	"golang.org/x/time/rate"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}
type Client struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

func perClientRateLimiter(next func(rw http.ResponseWriter, req *http.Request)) http.Handler {

	var (
		mu      sync.Mutex
		clients = make(map[string]*Client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.LastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ip, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &Client{
				Limiter: rate.NewLimiter(2, 4),
			}
		}
		clients[ip].LastSeen = time.Now()

		if !clients[ip].Limiter.Allow() {
			mu.Unlock()
			message := Message{
				Status: "Request Failed",
				Body:   "The API is at capacity, please try again later.",
			}
			rw.WriteHeader(http.StatusTooManyRequests) // 429
			json.NewEncoder(rw).Encode(&message)
			return
		}
		mu.Unlock()
		next(rw, req)
	})
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
	http.Handle("/ping", perClientRateLimiter(endpointHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("there was an error listening on port :8080", err)
	}
}
