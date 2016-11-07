package main

import (
	"./models"
	"./requesthandler"
	"net/http"
)

// NewEditHandler returns an instance of the edit handler, with
// the routes populated.
func NewEditHandler() *requesthandler.GenericRequestHandler {
	a := requesthandler.GenericRequestHandler{}
	a.RouteMap = map[string]requesthandler.Responder{
		"pages": pages,
	}
	return &a
}

func pages(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	return requesthandler.ResponseOK
}
