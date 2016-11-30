package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/colin353/markdown.ninja/models"
	"github.com/colin353/markdown.ninja/requesthandler"
)

// NewAuthenticationHandler creates an instance of the authentication
// handler and populates the routes hash.
func NewAuthenticationHandler() *requesthandler.GenericRequestHandler {
	a := requesthandler.GenericRequestHandler{}
	a.RouteMap = map[string]requesthandler.Responder{
		"login":        login,
		"check":        check,
		"logout":       logout,
		"signup":       signup,
		"check_domain": checkDomain,
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

// checkDomain looks takes a domain string and checks if it is available
// for registration or not.
func checkDomain(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	type checkDomainArgs struct {
		Domain string `json:"domain"`
	}
	args := checkDomainArgs{}
	err := requesthandler.ParseArguments(r, &args)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return requesthandler.ResponseInvalidArgs
	}

	user := models.User{}
	user.Domain = args.Domain
	err = models.Load(&user)
	if err != nil {
		return requesthandler.SimpleResponse{
			Result: "domain-available",
			Error:  false,
		}
	}

	return requesthandler.SimpleResponse{
		Result: "domain-exists",
		Error:  false,
	}
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

	// Try to create the user.
	me := models.NewUser()
	me.Name = args.Name
	me.Domain = args.Domain
	me.Email = args.Email

	// We don't really need an email address, except for when they need
	// to reset their password. So if they don't set an email address,
	// we'll just set their email to fake@fake.com
	if me.Email == "" {
		me.Email = "fake@fake.com"
	}

	me.SetPassword(args.Password)
	err = models.Insert(me)
	if err != nil {
		log.Printf("Failed to validate: %v", err.Error())
		return requesthandler.SimpleResponse{Result: "failed-validation", Error: true}
	}

	// All users will get a couple of files created for them
	// containing some basic defaults.
	defaultFiles, err := filepath.Glob("./web/default/*.md")
	for _, file := range defaultFiles {
		log.Printf("Creating default file: %s", file)
		p := models.Page{}
		p.Domain = me.Domain
		p.Name = filepath.Base(file)

		markdown, _ := ioutil.ReadFile(file)
		html, _ := ioutil.ReadFile(file + ".html")
		p.Markdown = string(markdown)
		p.HTML = string(html)
		models.Insert(&p)
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
