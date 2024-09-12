package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Message struct{
	Status string `json:"status"`
	Body string `json:"body"`
}

func perCLientRateLimiter(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	type client struct{
		limiter *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu sync.Mutex
		clients = make(map[string]*client)
	)

	// deleting the clients that have last seen of more than 3 minutes 
	go func ()  {
		for{
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients{
				if time.Since(client.lastSeen) > 3 *time.Minute{
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}	
	}()

	return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return 
		}

		mu.Lock()
		// if there is a new user 
		if _, found := clients[ip]; !found{
			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}

		// updating the last seen of client of that ip in the map 
		clients[ip].lastSeen = time.Now()


		// if the user is already in the map
		if !clients[ip].limiter.Allow(){
			mu.Unlock()
			w.WriteHeader(http.StatusTooManyRequests)
			message := Message{
				Status: "Request Failed",
				Body: "The API is at capacity, try again later",
			}
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(&message)
			return
		}
	mu.Unlock()
	// execute the next function
	next(w, r)
	})
}

func endpointHandler(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	message := Message{
		Status: "Successful",
		Body: "Hi!, You've reached the API. How may i help you?",
	}
	err :=json.NewEncoder(w).Encode(&message)
	if err != nil{
		return
	}
}
func main()  {
	http.Handle("/ping", perCLientRateLimiter(endpointHandler))
	err := http.ListenAndServe(":8000", nil)	
	if err != nil{
		log.Print("There was an error listening on port :8000", err)
	}
}