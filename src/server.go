package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/get", getPlaying)
	myRouter.HandleFunc("/{action}/{service}", controlMediaPlayer)

	myRouter.Use(loggingMiddleware)

	auth := os.Getenv("AUTH")
	if auth == "true" {
		authKey := os.Getenv("AUTH_KEY")
		if authKey == "" {
			authKey = "key123"
			log.Printf("No auth key found in .env, using `key123` as auth code.")
		}
		myRouter.Use(AuthMiddleware(authKey))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "10004"
	}

	serverAddress := ":" + port
	log.Printf("Server starting on %s", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, myRouter))
}

func respondWithError(w http.ResponseWriter, message string) {
	response := map[string]interface{}{
		"error":   true,
		"message": message,
	}
	jsonData, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonData)
}
