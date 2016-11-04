package main

import (
	_ "./models"
	"./requesthandler"
	"net/http"
	"github.com/gorilla/context"
)

// StaticAssetHandler serves static files, such as javascript,
// images, and HTML.
func StaticAssetHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" {
		http.ServeFile(w, r, "web/index.html")
	}
}

func main() {
	http.HandleFunc("/api/auth/", requesthandler.CreateHandler(NewAuthenticationHandler()))
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}
