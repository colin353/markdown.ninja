/*
  account.go

  Operations like changing your password, editing your account details,
  and setting up a custom domain name.
*/

package main

import (
	"log"
	"net/http"

	"github.com/colin353/markdown.ninja/models"
	"github.com/colin353/markdown.ninja/requesthandler"
)

// NewAccountHandler returns an instance of the account handler, with
// the routes populated.
func NewAccountHandler() *requesthandler.GenericRequestHandler {
	a := requesthandler.GenericRequestHandler{}
	a.RouteMap = map[string]requesthandler.Responder{
		"update_email":         updateEmail,
		"update_password":      updatePassword,
		"update_custom_domain": updateCustomDomain,
	}
	return &a
}

func updateEmail(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	type deleteArgs struct {
		Email string `json:"email"`
	}
	args := deleteArgs{}
	err := requesthandler.ParseArguments(r, &args)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return requesthandler.ResponseInvalidArgs
	}

	u.Email = args.Email
	err = models.Save(u)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return requesthandler.ResponseError
	}

	return requesthandler.ResponseOK
}

func updatePassword(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	type passwordArgs struct {
		Password string `json:"password"`
	}
	args := passwordArgs{}
	err := requesthandler.ParseArguments(r, &args)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return requesthandler.ResponseInvalidArgs
	}

	u.SetPassword(args.Password)
	err = models.Save(u)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return requesthandler.ResponseError
	}

	return requesthandler.ResponseOK
}

func updateCustomDomain(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	type domainArgs struct {
		Domain string `json:"domain"`
	}
	args := domainArgs{}
	err := requesthandler.ParseArguments(r, &args)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return requesthandler.ResponseInvalidArgs
	}

	log.Printf("got request for domain: %s", args.Domain)

	// If the user already has an external domain, we'll need to delete it.
	if u.ExternalDomain != "" {
		d := models.Domain{}
		d.ExternalDomain = u.ExternalDomain

		// We won't check for errors here, because if their
		// domain was not registered for some reason, that's
		// fine.
		models.Load(&d)

		// Need to check that the user actually owns
		// the domain in question.
		if d.InternalDomain == u.Domain {
			models.Delete(&d)
		}

		// Temporarily blank their external domain.
		u.ExternalDomain = ""
		err = models.Save(u)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return requesthandler.ResponseError
		}
	}

	domain := models.Domain{}
	domain.InternalDomain = u.Domain
	domain.ExternalDomain = args.Domain
	err = models.Insert(&domain)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return requesthandler.ResponseDuplicate
	}

	u.ExternalDomain = args.Domain
	err = models.Save(u)
	if err != nil {
		// Clean up by deleting the domain we created.
		models.Delete(&domain)

		http.Error(w, "", http.StatusBadRequest)
		return requesthandler.ResponseInvalidArgs
	}

	return requesthandler.ResponseOK
}
