/*
  prerender.go

  This file handles routing of the routes which should ultimately
  be handled by react.js. But to help make things a bit faster, we
  can do some prerendering. Cached prerendered versions of all of the
  relevant components are in ./web/prerendered/*.html, and a map of
  where to look up the prerendered stuff is in ./app/etc/routes.json.
*/

package requesthandler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Routes is a list of routes, which is loaded from the json file
// in ./app/etc/routes.json, and includes a list of prerendered
// components and how to serve them.
var Routes map[string]Route

// Route is an individual route, contianing a filename and a component.
type Route struct {
	File string `json:"file"`
}

func loadRoutes() {
	data, err := ioutil.ReadFile("./config/routes.json")
	if err != nil {
		panic("unable to load routes file ./config/routes.json")
	}
	err = json.Unmarshal(data, &Routes)
	if err != nil {
		panic("json could not be decoded in ./config/routes.json")
	}
}

// ReactHandler serves the index.html page to the group of
// subpages which are routed by react.js on the client.
func ReactHandler(w http.ResponseWriter, r *http.Request) {
	// If we have not yet loaded the routes map, let's do that now.
	if Routes == nil {
		loadRoutes()
	}

	// Check if the route we are rendering is available to be
	// prerendered.
	log.Printf("%v", r.RequestURI)
	route, ok := Routes[r.RequestURI]
	if ok {
		http.ServeFile(w, r, "web/prerendered/"+route.File)
	} else {
		http.ServeFile(w, r, "web/index.html")
	}
}
