package main

import (
	"./models"
	"./requesthandler"
	"github.com/gorilla/context"
	"log"
	"net/http"
)

// indexHandler serves the index.html page to the group of
// subpages which are routed by react.js on the client.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func main() {
	// Set up routing.
	http.HandleFunc("/api/auth/", requesthandler.CreateHandler(NewAuthenticationHandler()))
	http.HandleFunc("/api/edit/", requesthandler.CreateAuthenticatedHandler(NewEditHandler()))
	http.HandleFunc("/edit/", indexHandler)
	http.Handle("/", http.FileServer(http.Dir("./web")))

	// Connect to redis.
	models.Connect()

	// Start up the server.
	err := http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatalf("Unable to start server: %v", err.Error())
	}
}
