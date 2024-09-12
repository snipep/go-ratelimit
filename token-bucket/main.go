package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Message struct{
	Status string `json:"status"`
	Body string `json:"body"` 
}

func endpointHandler(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	message := Message{
		Status: "Successful",
		Body: "Hi!, You're reached the API. HOw may i help you?",
	}
	err := json.NewEncoder(w).Encode(&message)
	if err != nil{
		return
	}
}
func main()  {
	http.Handle("/ping", rateLimiter(endpointHandler))
	err := http.ListenAndServe(":8000", nil)
	if err != nil{
		log.Print("There was an error listening on port :8080)", err)
	}
}