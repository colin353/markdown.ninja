package main

import (
	"github.com/colin353/portfolio/models"
	"github.com/colin353/portfolio/requesthandler"
	"log"
	"net/http"
)

// NewAuthenticationHandler creates an instance of the authentication
// handler and populates the routes hash.
func NewAuthenticationHandler() *requesthandler.GenericRequestHandler {
	a := requesthandler.GenericRequestHandler{}
	a.RouteMap = map[string]requesthandler.Responder{
		"login":  login,
		"check":  check,
		"logout": logout,
		"signup": signup,
	}
	return &a
}

// The check function tries to figure out if we are currently logged in. The
// react installation requires it, because it needs to make routing decisions based
// upon the authentication state, but ultimately it is up to the server, not the client
// to decide who is authenticated.
func check(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	authenticated, user := requesthandler.CheckAuthentication(w, r)

	if authenticated {
		return map[string]interface{}{
			"status": "ok",
			"error":  false,
			"user":   user.Export(),
		}
	}
	http.Error(w, "Not authorized", http.StatusForbidden)
	return requesthandler.ResponseError
}

func login(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	type loginArgs struct {
		Domain   string `json:"domain"`
		Password string `json:"password"`
	}
	args := loginArgs{}
	err := requesthandler.ParseArguments(r, &args)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return requesthandler.ResponseInvalidArgs
	}

	// Check if the user is in the database, and if their password matches.
	me := models.User{}
	me.Domain = args.Domain
	err = models.Load(&me)
	if err != nil {
		// The user doesn't exist in the database.
		log.Printf("User @ domain `%v` doesn't exist.\n", args.Domain)
		return requesthandler.ResponseError
	}

	// Check the password.
	if !me.CheckPassword(args.Password) {
		log.Printf("User @ domain `%v`: wrong password.\n", args.Domain)
		return requesthandler.ResponseError
	}

	// The user has met the authentication requirements, so we will write
	// their cookie.
	session, _ := requesthandler.SessionStore.Get(r, "authentication")
	session.Values["authenticated"] = true
	session.Values["domain"] = me.Domain
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Failed to save session.")
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return requesthandler.ResponseError
	}

	return requesthandler.ResponseOK
}

func logout(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	session, _ := requesthandler.SessionStore.Get(r, "authentication")
	// The user has met the authentication requirements, so we will write
	// their cookie.
	session.Options.MaxAge = -1
	err := session.Save(r, w)
	if err != nil {
		log.Printf("Failed to save session.")
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return requesthandler.ResponseError
	}

	return requesthandler.ResponseOK
}

func signup(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	type signupArgs struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Domain   string `json:"domain"`
		Password string `json:"password"`
	}
	args := signupArgs{}
	err := requesthandler.ParseArguments(r, &args)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return requesthandler.ResponseInvalidArgs
	}

	// Check the signup conditions.
	if len(args.Password) < 6 {
		return requesthandler.SimpleResponse{
			Result: "password-too-short",
			Error:  true,
		}
	}

	me := models.NewUser()
	me.Name = args.Name
	me.Domain = args.Domain
	me.Email = args.Email
	me.SetPassword(args.Password)
	err = models.Insert(me)
	if err != nil {
		log.Printf("Failed to validate: %v", err.Error())
		return requesthandler.SimpleResponse{"failed-validation", true}
	}

	// I guess we created the user OK, so let's log them in also.
	session, _ := requesthandler.SessionStore.Get(r, "authentication")
	session.Values["authenticated"] = true
	session.Values["domain"] = me.Domain
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Failed to save session.")
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return requesthandler.ResponseError
	}

	return requesthandler.ResponseOK
}
