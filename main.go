package main

import (
	"log"
	"net/http"

	"github.com/colin353/markdown.ninja/config"
	"github.com/colin353/markdown.ninja/models"
	"github.com/colin353/markdown.ninja/requesthandler"
	"github.com/gorilla/context"
)

// AppConfig contains the application configuration, which is loaded
// from the config yaml files and overridden by environment variables.
var AppConfig *config.Config

func main() {
	// Load the configuration file, and distribute it to the modules.
	AppConfig = config.LoadConfig("./config")
	models.AppConfig = AppConfig
	requesthandler.AppConfig = AppConfig

	// Connect to redis.
	models.Connect()

	// If we are in testing mode, we must delete the database contents.
	if AppConfig.Mode == "test" || AppConfig.Mode == "testing" {
		models.ClearDatabase()
	}

	// Set up routing.
	http.HandleFunc("/api/auth/", requesthandler.CreateHandler(NewAuthenticationHandler()))
	http.HandleFunc("/api/edit/", requesthandler.CreateAuthenticatedHandler(NewEditHandler()))
	http.HandleFunc("/api/files/", requesthandler.CreateAuthenticatedHandler(NewFileHandler()))
	http.HandleFunc("/edit/", requesthandler.ReactHandler)

	http.HandleFunc("/", requesthandler.SubdomainHandler)

	// Start up the server.
	err := http.ListenAndServe(":"+AppConfig.Port, context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatalf("Unable to start server: %v", err.Error())
	}
}
