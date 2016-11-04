package requesthandler

import (
  "log"
  "net/http"
  "strings"
  "github.com/gorilla/sessions"
)

// SessionStore saves cookie data securely.
var SessionStore = sessions.NewCookieStore([]byte("123test"))

// A Responder is a function which can respond directly to an HTTP
// request.
type Responder func(http.ResponseWriter, *http.Request)

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
func CreateHandler(rh RequestHandler) Responder {
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
      responder(w, r)
    }
  }
}

// CreateAuthenticatedHandler takes a RequestHandler and wraps it with
// authentication requirements so that the endpoints cannot be accessed
// without logging in first.
func CreateAuthenticatedHandler(rh RequestHandler) Responder {
  return func(w http.ResponseWriter, r *http.Request) {
    // Check if the authentication requirements are met
    session, err := SessionStore.Get(r, "authentication")
    if err != nil {
      log.Printf("Could not open session store.",)
      http.Error(w, "Internal server error.", http.StatusInternalServerError)
    }

    authenticated, ok := session.Values["authenticated"].(bool)
    if !ok {
      authenticated = false
    }
    if !authenticated {
      log.Printf("401: not authorized to access `%v`", r.URL.Path)
      http.Error(w, "Not authorized", http.StatusForbidden)
    } else {
      // We are authenticated, so just execute the normal handler.
      CreateHandler(rh)(w, r)
    }
  }
}
