package main

import (
	"./requesthandler"
	"log"
	"net/http"
)

// NewAuthenticationHandler creates an instance of the authentication
// handler and populates the routes hash.
func NewAuthenticationHandler() *requesthandler.GenericRequestHandler {
	a := requesthandler.GenericRequestHandler{}
	a.RouteMap = map[string]requesthandler.Responder{
		"login": login,
	}
	return &a
}

func login(w http.ResponseWriter, r *http.Request) {
	session, err := requesthandler.SessionStore.Get(r, "authentication")
	if err != nil {
		log.Printf("Could not open session store.")
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
	}

	// The user has met the authentication requirements, so we will write
	// their cookie.
	session.Values["authenticated"] = true
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Failed to save session.")
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
	}

	w.Write([]byte("You have logged yourself in."))
}
