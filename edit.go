package main

import (
	"log"
	"net/http"

	"github.com/colin353/portfolio/models"
	"github.com/colin353/portfolio/requesthandler"
)

// NewEditHandler returns an instance of the edit handler, with
// the routes populated.
func NewEditHandler() *requesthandler.GenericRequestHandler {
	a := requesthandler.GenericRequestHandler{}
	a.RouteMap = map[string]requesthandler.Responder{
		"pages":       pages,
		"create_page": createPage,
		"edit_page":   editPage,
	}
	return &a
}

// Create a new page. If that page already exists, will return an error.
func createPage(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	type createArgs struct {
		Markdown string `json:"markdown"`
		HTML     string `json:"html"`
	}
	args := createArgs{}
	err := requesthandler.ParseArguments(r, &args)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return requesthandler.ResponseInvalidArgs
	}

	// Create a new instance of the page object.
	p := models.Page{
		Markdown: args.Markdown,
		HTML:     args.HTML,
		Domain:   u.Domain,
	}
	err = p.GenerateName()
	if err != nil {
		log.Printf("Tried to create a new page, but couldn't make a unique name. (tried %s)", p.Key())
		http.Error(w, "", http.StatusBadRequest)
	}

	err = models.Insert(&p)
	if err != nil {
		log.Printf("Tried to create a new page called `%s`, but encountered an error.", p.Key())
		http.Error(w, "", http.StatusBadRequest)
	}

	return p.Export()
}

func editPage(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	type editArgs struct {
		Name     string `json:"name"`
		Markdown string `json:"markdown"`
		HTML     string `json:"html"`
	}
	args := editArgs{}
	err := requesthandler.ParseArguments(r, &args)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return requesthandler.ResponseInvalidArgs
	}

	// Load the old version of the page.
	p := models.Page{}
	p.Domain = u.Domain
	p.Name = args.Name
	err = models.Load(&p)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return requesthandler.ResponseInvalidArgs
	}

	// Update the data and save it.
	p.Markdown = args.Markdown
	p.HTML = args.HTML
	err = models.Save(&p)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return requesthandler.ResponseInvalidArgs
	}

	return requesthandler.ResponseOK
}

// Return a list of pages belonging to that user.
func pages(u *models.User, w http.ResponseWriter, r *http.Request) interface{} {
	// Create a page object and use it to search for its own siblings.
	p := models.Page{}
	p.Domain = u.Domain
	iterator, err := models.GetList(&p)
	if err != nil {
		log.Printf("Tried to load pages under `%s`, but it failed.", p.RegistrationKey())
		return requesthandler.ResponseError
	}

	pageList := make([]map[string]interface{}, 0, iterator.Count())
	for iterator.Next() {
		pageList = append(pageList, iterator.Value().Export())
	}

	return pageList
}
