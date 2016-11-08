/*
  subdomain.go

  This file handles rendering a page visited from a subdomain.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/colin353/portfolio/models"
)

func renderSubdomain(domain string, w http.ResponseWriter, r *http.Request) {
	p := models.Page{}
	p.Domain = domain
	p.Name = "index.md"
	err := models.Load(&p)
	if err != nil {
		log.Printf("Didn't find anything at key: `%v`", p.Key())
		http.Error(w, "404: that thing doesn't exist!", http.StatusNotFound)
		return
	}

	defaultStyle, err := ioutil.ReadFile("web/css/webstyles/default.css")
	if err != nil {
		log.Println("Could not open required style file: web/css/webstyles/default.css")
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		return
	}

	requiredStyle, err := ioutil.ReadFile("web/css/webstyles/required.css")
	if err != nil {
		log.Println("Could not open required style file: web/css/webstyles/default.css")
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf(`
      <style>
        %s
        %s
      </style>
      <div class='container'>
        <div class='content'>%s</div>
      </div>
    `, requiredStyle, defaultStyle, p.HTML)))
}
