package main

import (
	"log"
	"net/http"

	"github.com/colin353/portfolio/config"
	"github.com/colin353/portfolio/models"
	"github.com/colin353/portfolio/requesthandler"
	"github.com/gorilla/context"
)

// AppConfig contains the application configuration, which is loaded
// from the config yaml files and overridden by environment variables.
var AppConfig *config.Config

// indexHandler serves the index.html page to the group of
// subpages which are routed by react.js on the client.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func main() {
	// Load the configuration file, and distribute it to the modules.
	AppConfig = config.LoadConfig("./config")
	models.AppConfig = AppConfig
	requesthandler.AppConfig = AppConfig

	// Set up routing.
	http.HandleFunc("/api/auth/", requesthandler.CreateHandler(NewAuthenticationHandler()))
	http.HandleFunc("/api/edit/", requesthandler.CreateAuthenticatedHandler(NewEditHandler()))
	http.HandleFunc("/api/files/", requesthandler.CreateAuthenticatedHandler(NewFileHandler()))
	http.HandleFunc("/edit/", indexHandler)
	//http.HandleFunc("/domain", handleSubdomain)
	http.HandleFunc("/", requesthandler.SubdomainHandler)

	// Connect to redis.
	models.Connect()

	// Start up the server.
	err := http.ListenAndServe(":"+AppConfig.Port, context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatalf("Unable to start server: %v", err.Error())
	}
}
