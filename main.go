package main

import (
  "net/http"
  _ "./models"
)

func handler(w http.ResponseWriter, r *http.Request) {
  if(r.URL.Path == "") {
    http.ServeFile(w, r, "web/index.html")
  }
}

func main() {
  http.HandleFunc("/", handler)
  //http.Handle("/", http.FileServer(http.Dir("./web")))
  http.ListenAndServe(":8080", nil)
}
