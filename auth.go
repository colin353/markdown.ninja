package main

import (
	"./requesthandler"
	"./models"
	"log"
	"net/http"
)

// NewAuthenticationHandler creates an instance of the authentication
// handler and populates the routes hash.
func NewAuthenticationHandler() *requesthandler.GenericRequestHandler {
	a := requesthandler.GenericRequestHandler{}
	a.RouteMap = map[string]requesthandler.Responder{
		"login": login,
		"check": check,
		"logout": logout,
	}
	return &a
}

// The check function tries to figure out if we are currently logged in. The
// react installation requires it, because it needs to make routing decisions based
// upon the authentication state, but ultimately it is up to the server, not the client
// to decide who is authenticated.
func check(u *models.User, args interface{}, w http.ResponseWriter, r *http.Request) {
	authenticated := requesthandler.CheckAuthentication(w, r)

	if authenticated {
		w.Write([]byte("true"))
	} else {
		http.Error(w, "Not authorized", http.StatusForbidden)
	}
}

func login(u *models.User, args interface{}, w http.ResponseWriter, r *http.Request) {
	session, _ := requesthandler.SessionStore.Get(r, "authentication")
	// The user has met the authentication requirements, so we will write
	// their cookie.
	session.Values["authenticated"] = true
	err := session.Save(r, w)
	if err != nil {
		log.Printf("Failed to save session.")
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
	}

	w.Write([]byte("You are logged in."))
}

func logout(u *models.User, args interface{}, w http.ResponseWriter, r *http.Request) {
	session, _ := requesthandler.SessionStore.Get(r, "authentication")
	// The user has met the authentication requirements, so we will write
	// their cookie.
	session.Options.MaxAge = -1
	err := session.Save(r, w)
	if err != nil {
		log.Printf("Failed to save session.")
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
	}

	w.Write([]byte("You are logged out."))
}
