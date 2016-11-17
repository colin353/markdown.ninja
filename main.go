package main

import (
	"log"
	"net/http"
	"strings"

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

// This handler determines whether to serve subdomain content or not. If
// it determines that the request is for a subdomain, it'll hand off that
// request to the subdomain renderer.
func subdomainHandler(w http.ResponseWriter, r *http.Request) {
	// In case we are running this from an IP address, or from a subdomain,
	// we need to make sure we're looking at the subdomain lower than the
	// one we recognize as our own.
	domains := strings.Split(r.Host, ".")
	if len(domains) == 1 {
		http.FileServer(http.Dir("./web")).ServeHTTP(w, r)
		return
	}
	renderSubdomain(domains[0], w, r)
}

func main() {
	// Load the configuration file, and distribute it to the modules.
	AppConfig = config.LoadConfig("./config")
	models.AppConfig = AppConfig

	// Set up routing.
	http.HandleFunc("/api/auth/", requesthandler.CreateHandler(NewAuthenticationHandler()))
	http.HandleFunc("/api/edit/", requesthandler.CreateAuthenticatedHandler(NewEditHandler()))
	http.HandleFunc("/api/files/", requesthandler.CreateAuthenticatedHandler(NewFileHandler()))
	http.HandleFunc("/edit/", indexHandler)
	//http.HandleFunc("/domain", handleSubdomain)
	http.HandleFunc("/", subdomainHandler)

	// Connect to redis.
	models.Connect()

	// Start up the server.
	err := http.ListenAndServe(":"+AppConfig.Port, context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatalf("Unable to start server: %v", err.Error())
	}
}
