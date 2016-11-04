package requesthandler

import (
	"../models"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"strings"
)

// SessionStore saves cookie data securely.
var SessionStore = sessions.NewCookieStore([]byte("123t22est"))

// A Responder is a function which can respond directly to an HTTP
// request.
type Responder func(*models.User, interface{}, http.ResponseWriter, *http.Request)

// An IntermediateResponder is what must respond first to an HTTP
// request. It's the normal type required for a handler.
type IntermediateResponder func(http.ResponseWriter, *http.Request)

// A RequestHandler represents a collection of responders which all
// fit into the same category.
type RequestHandler interface {
	Route(string) Responder
}

// GenericRequestHandler contains a routemap and just calls one of the
// route functions when the path is satisfiied.
type GenericRequestHandler struct {
	RouteMap map[string]Responder
}

// Route satisfies the interface requirements in the requesthandler module
// so that the AuthenticationHandler can be used to service requests
func (rh *GenericRequestHandler) Route(route string) Responder {
	responder, ok := rh.RouteMap[route]
	if !ok {
		return nil
	}
	return responder
}

// CreateHandler takes a RequestHandler and turns it into a function
// which can respond to HTTP requests by returning an anonymous function
// bound with the RequestHandler.
func CreateHandler(rh RequestHandler) IntermediateResponder {
	return func(w http.ResponseWriter, r *http.Request) {
		// We'll need to break down the URL path to get the correct routing.
		paths := strings.Split(r.URL.RequestURI()[1:], "/")

		// Routing example:
		// --> /api/auth/login
		// The first part of the request URI is "api". All requesthandler
		// endpoints use /api, so this is not used for routing. The second part
		// of the request uri distinguishes different handlers, so by now it has
		// already been accounted for. So we need to use the third piece to route on.
		var responder Responder
		if len(paths) == 3 {
			responder = rh.Route(paths[2])
		}
		if responder == nil {
			// Must have been a 404.
			log.Printf("404: no such path `%v`", r.URL.Path)
			http.Error(w, "No such path", http.StatusNotFound)
		} else {
			responder(nil, nil, w, r)
		}
	}
}

// CheckAuthentication uses the current session and request variables to check
// if the authentication requirements are met. It returns true if met. Might raise
// an error if something goes wrong: but it automatically reports status 500, so
// no need to handle.
func CheckAuthentication(w http.ResponseWriter, r *http.Request) bool {
	// If an error occurs while loading the session, it may be because the client
	// provided invalid information so we will just report them as an illegal login.
	session, err := SessionStore.Get(r, "authentication")
	if err != nil {
		return false
	}

	authenticated, ok := session.Values["authenticated"].(bool)
	if !ok {
		return false
	}
	return authenticated
}

// CreateAuthenticatedHandler takes a RequestHandler and wraps it with
// authentication requirements so that the endpoints cannot be accessed
// without logging in first.
func CreateAuthenticatedHandler(rh RequestHandler) IntermediateResponder {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated := CheckAuthentication(w, r)

		// Check if the authentication requirements are met
		if !authenticated {
			log.Printf("401: not authorized to access `%v`", r.URL.Path)
			http.Error(w, "Not authorized", http.StatusForbidden)
		} else {
			// We are authenticated, so just execute the normal handler.
			CreateHandler(rh)(w, r)
		}
	}
}
