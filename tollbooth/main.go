package main

import (
	"encoding/json"
	"log"
	"net/http"

	tollbooth "github.com/didip/tollbooth/v7"
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
	message := Message{
		Status: "Request Failed",
		Body: "The API is at capacity, try again later",
	}
	jsonMessage, _ := json.Marshal(message)
	tbthLimiter := tollbooth.NewLimiter(2, nil)
	tbthLimiter.SetMessageContentType("application/json")
	tbthLimiter.SetMessage(string(jsonMessage))
	http.Handle("/ping", tollbooth.LimitFuncHandler(tbthLimiter, endpointHandler))
	err := http.ListenAndServe(":8000", nil)
	if err != nil{
		log.Print("There was an error listening on port :8000", err)
	}
}	