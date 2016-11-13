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

// There are really only two possible routes for a subdomain. We will either
// be serving a Page, which is basically some HTML wrapped in some elements and
// supported by some CSS, or we will be serving a file.
//
// Files are served from any path starting with /files/*, and pages are served
// for everything else. This function checks the request string and figure out
// whether we should be responding with a page or a file, and calls the appropriate
// function.
func renderSubdomain(domain string, w http.ResponseWriter, r *http.Request) {
	if len(r.RequestURI) < 7 || r.RequestURI[0:7] != "/files/" {
		renderPage(domain, w, r)
	} else {
		renderFile(domain, w, r)
	}
}

func renderPage(domain string, w http.ResponseWriter, r *http.Request) {
	p := models.Page{}
	p.Domain = domain
	log.Printf("Request URI: %s", r.RequestURI)

	if r.RequestURI == "/" || r.RequestURI == "/index.md" {
		p.Name = "index.md"
	} else {
		p.Name = r.RequestURI[1:]
	}

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

func renderFile(domain string, w http.ResponseWriter, r *http.Request) {
	f := models.File{}
	f.Domain = domain
	f.Name = r.RequestURI[7:]
	err := models.Load(&f)

	if err != nil {
		log.Printf("Didn't find anything at key: `%v`", f.Key())
		http.Error(w, "404: that thing doesn't exist!", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, f.GetPath())
}
