package main

import (
	"encoding/json"
	"golang.org/x/time/rate"
	"net/http"
)

// rateLimiter is a middleware that limits the number of requests per second
func rateLimiter(next func(wr http.ResponseWriter, req *http.Request)) http.HandlerFunc {

	// Create a limiter that allows 2 requests per second with a maximum burst of 4
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		// If the limiter is not allowing the request, return an error
		if !limiter.Allow() {
			message := Message{
				Status: "Request Failed",
				Body:   "The API is at capacity, please try again later.",
			}
			rw.WriteHeader(http.StatusTooManyRequests) // 429
			json.NewEncoder(rw).Encode(&message)
			return
		} else {
			// If the limiter allows the request, call the next handler
			next(rw, req)
		}
	})
}
