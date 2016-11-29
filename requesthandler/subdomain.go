/*
  subdomain.go

  This file handles rendering a page visited from a subdomain.
*/

package requesthandler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/colin353/portfolio/models"
)

// SubdomainHandler determines whether to serve subdomain content or not. If
// it determines that the request is for a subdomain, it'll hand off that
// request to the subdomain renderer.
func SubdomainHandler(w http.ResponseWriter, r *http.Request) {
	// In case we are running this from an IP address, or from a subdomain,
	// we need to make sure we're looking at the subdomain lower than the
	// one we recognize as our own.
	subdomain := getSubdomainFromHost(r.Host, AppConfig.Hostnames)
	if subdomain == "" {
		// Special case: if we're actually looking at the index page,
		// we need to consider using the prerendered files.
		if r.RequestURI == "/" {
			log.Println("Served prerendered index.")
			http.ServeFile(w, r, "web/prerendered/index.html")
		} else {
			http.FileServer(http.Dir("./web")).ServeHTTP(w, r)
		}
		return
	}
	renderSubdomain(subdomain, w, r)
}

// This function takes a host string, like you might get from a
// HTTP request (e.g. something like "subdomain.hostname.hostname.com:3060")
// and figures out which components of that subdomain are actually the subdomain.
// In the example above it will return "subdomain". If no subdomain is found,
// it returns an empty string.
func getSubdomainFromHost(host string, hostnames []string) string {
	// Strip the port number, if it is set.
	port := strings.Index(host, ":")
	if port != -1 {
		host = host[:port]
	}

	// Check if we can find our own hostname, and if so, delete it.
	for _, h := range hostnames {
		if strings.HasSuffix(host, h) {
			host = host[:len(host)-len(h)]
			break
		}
	}

	// Split the resulting string into an array of dot-separated
	// subdomain strings. Remove any empty components, and return
	// the result.
	parts := []string{}
	for _, part := range strings.Split(host, ".") {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return strings.Join(parts, ".")
}

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

	user := models.User{}
	user.Domain = domain
	models.Load(&user)

	defaultStyle, err := ioutil.ReadFile(fmt.Sprintf("web/css/webstyles/%s.css", user.Style))
	if err != nil {
		log.Println("Could not open required style file: web/css/webstyles/default.css")
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		return
	}

	requiredStyle, err := ioutil.ReadFile("web/css/webstyles/required.css")
	if err != nil {
		log.Println("Could not open required style file: web/css/webstyles/required.css")
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf(`
      <style>
        %s
        %s
      </style>
      <div class='md_container'>
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
