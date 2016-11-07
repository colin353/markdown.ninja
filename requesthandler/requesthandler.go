package requesthandler

import (
	"github.com/colin353/portfolio/models"
	"encoding/json"
	"github.com/gorilla/sessions"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// A SimpleResponse is just an acknowledgement response
// which reports if an error occurred or not.
type SimpleResponse struct {
	Result string `json:"result"`
	Error  bool   `json:"error"`
}

// These are the response types that should be returned by the
// Responders. They're just JSON strings.
var (
	ResponseOK          = SimpleResponse{"ok", false}
	ResponseError       = SimpleResponse{"error", true}
	ResponseInvalidArgs = SimpleResponse{"invalid arguments", true}
)

// SessionStore saves cookie data securely.
var SessionStore = sessions.NewCookieStore([]byte("123t22est"))

// A Responder is a function which can respond directly to an HTTP
// request.
type Responder func(*models.User, http.ResponseWriter, *http.Request) interface{}

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

// ParseArguments takes a struct of the desired type and tries to convert
// the provided arguments
func ParseArguments(r *http.Request, args interface{}) error {
	formData, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err
	}
	return json.Unmarshal(formData, args)
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
		// Since there's no user (they may not be authenticated)
		// the user object is just nil.
		routeRequest(rh, nil, w, r)
	}
}

func routeRequest(rh RequestHandler, u *models.User, w http.ResponseWriter, r *http.Request) {
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
		response := responder(nil, w, r)

		// The response should be an object that can be converted to JSON.
		// If it is nil, we will assume that the result is OK.
		if response == nil {
			response = &ResponseOK
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("Failed to Marshal response `%v` to JSON", response)
			http.Error(w, "Internal error", http.StatusInternalServerError)
		}

		w.Write(responseJSON)
	}
}

// CheckAuthentication uses the current session and request variables to check
// if the authentication requirements are met. It returns true if met. Might raise
// an error if something goes wrong: but it automatically reports status 500, so
// no need to handle.
func CheckAuthentication(w http.ResponseWriter, r *http.Request) (bool, *models.User) {
	// If an error occurs while loading the session, it may be because the client
	// provided invalid information so we will just report them as an illegal login.
	session, err := SessionStore.Get(r, "authentication")
	if err != nil {
		return false, nil
	}

	authenticated, ok := session.Values["authenticated"].(bool)
	if !ok {
		return false, nil
	}

	domain, ok := session.Values["domain"].(string)
	if !ok {
		return false, nil
	}

	// It's also necessary to check that the user record in the datbaase is valid.
	user := models.User{}
	user.Domain = domain
	err = models.Load(&user)
	if err != nil {
		// The record doesn't exist: so they are not authenticatd. In addition to
		// returning false, we'll also delete their invalid cookie.
		session.Options.MaxAge = -1
		session.Save(r, w)

		return false, nil
	}

	return authenticated, &user
}

// CreateAuthenticatedHandler takes a RequestHandler and wraps it with
// authentication requirements so that the endpoints cannot be accessed
// without logging in first.
func CreateAuthenticatedHandler(rh RequestHandler) IntermediateResponder {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticated, user := CheckAuthentication(w, r)

		// Check if the authentication requirements are met
		if !authenticated {
			log.Printf("401: not authorized to access `%v`", r.URL.Path)
			http.Error(w, "Not authorized", http.StatusForbidden)
			return
		}

		// We are authenticated, so just execute the normal handler.
		routeRequest(rh, user, w, r)
	}
}
